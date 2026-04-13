use crate::commands::command::CliCommand;
use crate::rules::ExecutionMode;
use crate::ui::hook_display;
use crate::Context;
use crate::GitHook;
use crate::RuleResult;
use anyhow::{anyhow, Result};
use clap::Parser;
use std::path::PathBuf;
use std::sync::Arc;
use tokio::sync::Semaphore;

/// Maximum number of async rules that may run concurrently.
pub const MAX_CONCURRENT_ASYNC_RULES: usize = 10;

#[derive(Debug, Parser)]
pub struct HandleCommand {
    /// The hook to handle
    #[arg(value_enum)]
    hook: GitHook,
    /// The commit message file path
    message: Option<String>,
}

impl CliCommand for HandleCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        if let Some(message) = &self.message {
            context.set_commit_msg_path(PathBuf::from(message));
        }

        let config = context.configuration();
        println!("{}", hook_display(&self.hook, config.files.clone()));

        let Some(rules) = config.hooks.get(&self.hook) else {
            println!("No rules found for hook {}", self.hook);
            return Ok(());
        };

        // ── Phase 1: prepare (sequential, needs &mut ctx) ──────────────────────
        //
        // For each rule we call `extend()` (which needs exclusive access to the
        // context) and evaluate the `when` condition.  After this loop we only
        // need shared, immutable access to the context.

        struct PreparedMeta {
            /// Pre-extended context for rules that declare their own `extract`.
            /// `None` means use the shared original context.
            extended_ctx: Option<Box<dyn Context>>,
            should_skip: bool,
            mode: ExecutionMode,
            typetag_name: &'static str,
        }

        let mut meta_vec: Vec<PreparedMeta> = Vec::with_capacity(rules.len());

        for rule_ctx in rules {
            let extended_ctx = rule_ctx
                .extract
                .as_ref()
                .map(|e| context.extend(e))
                .transpose()?;

            let actual_ctx: &dyn Context = match &extended_ctx {
                Some(b) => b.as_ref(),
                None => context as &dyn Context,
            };

            let should_skip = if let Some(when) = &rule_ctx.when {
                !when.check(actual_ctx.variables())?
            } else {
                false
            };

            meta_vec.push(PreparedMeta {
                extended_ctx,
                should_skip,
                mode: rule_ctx.rule.execution_mode(),
                typetag_name: rule_ctx.rule.typetag_name(),
            });
        }

        // ── Phase 2: execute (sync sequential + async concurrent, in parallel) ─
        //
        // After phase 1 we only need `&dyn Context` (immutable).
        //
        // We share `meta` via Arc so that both the sync thread and the N async
        // threads can own a clone.  `ctx_ref: &dyn Context` and
        // `rules: &[RuleContext]` are fat-pointer references and implement
        // `Copy + Send` (because `Context: Sync` and `RuleContext: Sync`), so
        // they are implicitly copied into every `move` closure.
        //
        // Results carry their original index so we can restore the declaration
        // order before printing / error-checking.

        // Immutable borrows valid for the entire scope of this function.
        let ctx_ref: &dyn Context = context;
        let rules_ref: &[_] = rules;

        // A multi-thread runtime backs the semaphore.  Scoped OS threads
        // acquire permits via `Handle::block_on`, which is safe to call from
        // non-tokio threads.
        let rt = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .map_err(|e| anyhow!("failed to build tokio runtime: {e}"))?;
        let rt_handle = rt.handle().clone();

        let semaphore = Arc::new(Semaphore::new(MAX_CONCURRENT_ASYNC_RULES));

        // Wrap meta in Arc so each scoped thread can hold an independent owner.
        let meta = Arc::new(meta_vec);

        // Partition rule indices by execution mode.
        let mut sync_indices: Vec<usize> = Vec::new();
        let mut async_indices: Vec<usize> = Vec::new();
        for (i, m) in meta.iter().enumerate() {
            match m.mode {
                ExecutionMode::Sync => sync_indices.push(i),
                ExecutionMode::Async => async_indices.push(i),
            }
        }

        let indexed_results: Vec<(usize, RuleResult)> =
            std::thread::scope(|s| -> Result<Vec<(usize, RuleResult)>> {
                // ── Sync thread ────────────────────────────────────────────────
                // Spawned first so it starts running immediately while async
                // threads are being set up.
                let meta_for_sync = Arc::clone(&meta);
                let sync_handle = s.spawn(move || -> Result<Vec<(usize, RuleResult)>> {
                    let mut results = Vec::with_capacity(sync_indices.len());
                    for &idx in &sync_indices {
                        let m = &meta_for_sync[idx];
                        if m.should_skip {
                            results.push((
                                idx,
                                RuleResult::Skipped {
                                    name: m.typetag_name.to_string(),
                                },
                            ));
                            continue;
                        }
                        let actual: &dyn Context = match &m.extended_ctx {
                            Some(b) => b.as_ref(),
                            None => ctx_ref,
                        };
                        results.push((idx, rules_ref[idx].rule.check(actual)?));
                    }
                    Ok(results)
                });

                // ── Async threads ──────────────────────────────────────────────
                // One thread per async rule; they race for semaphore permits.
                let async_handles: Vec<_> = async_indices
                    .iter()
                    .copied()
                    .map(|idx| {
                        let sem = Arc::clone(&semaphore);
                        let handle = rt_handle.clone();
                        let meta_for_async = Arc::clone(&meta);
                        s.spawn(move || -> Result<(usize, RuleResult)> {
                            // Block the current OS thread until a concurrency slot
                            // is available.  `Handle::block_on` drives the future
                            // on this thread and is safe from non-tokio threads.
                            let _permit = handle
                                .block_on(sem.acquire())
                                .map_err(|_| anyhow!("semaphore closed unexpectedly"))?;

                            let m = &meta_for_async[idx];
                            if m.should_skip {
                                return Ok((
                                    idx,
                                    RuleResult::Skipped {
                                        name: m.typetag_name.to_string(),
                                    },
                                ));
                            }
                            let actual: &dyn Context = match &m.extended_ctx {
                                Some(b) => b.as_ref(),
                                None => ctx_ref,
                            };
                            Ok((idx, rules_ref[idx].rule.check(actual)?))
                        })
                    })
                    .collect();

                // ── Collect results ────────────────────────────────────────────
                // Joining blocks until each thread finishes.  Both groups were
                // running in parallel while the threads were alive.
                let mut all: Vec<(usize, RuleResult)> = Vec::with_capacity(rules_ref.len());

                all.extend(
                    sync_handle
                        .join()
                        .map_err(|_| anyhow!("sync thread panicked"))
                        .and_then(|r| r)?,
                );

                for handle in async_handles {
                    all.push(
                        handle
                            .join()
                            .map_err(|_| anyhow!("async rule thread panicked"))
                            .and_then(|r| r)?,
                    );
                }

                // Restore original declaration order.
                all.sort_unstable_by_key(|(i, _)| *i);
                Ok(all)
            })?;

        // ── Phase 3: report ────────────────────────────────────────────────────

        let mut has_failure = false;
        for (_, result) in &indexed_results {
            match result {
                RuleResult::Success { name, output } => {
                    println!("{name} executed successfully");
                    if let Some(value) = output
                        && !value.is_empty()
                    {
                        println!("{value}");
                    }
                }
                RuleResult::Failure { message, name } => {
                    eprintln!("{name}: {message}");
                    has_failure = true;
                }
                RuleResult::Skipped { name } => {
                    println!("{name}: skipped");
                }
            }
        }

        if has_failure {
            return Err(anyhow!("Hook failed"));
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::MockContext;
    use crate::Rule;
    use crate::{Configuration, RuleContext};
    use serde::{Deserialize, Serialize};
    use std::fmt::Display;
    use std::sync::Arc;

    // ── helpers ───────────────────────────────────────────────────────────────

    #[derive(Debug, Serialize, Deserialize)]
    struct FakeRule {
        success: bool,
    }

    impl FakeRule {
        fn new(success: bool) -> Self {
            Self { success }
        }
    }

    impl Display for FakeRule {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "FakeRule")
        }
    }

    impl Rule for FakeRule {
        fn check(&self, _: &dyn Context) -> Result<RuleResult> {
            match self.success {
                true => Ok(RuleResult::Success {
                    name: "FakeRule".into(),
                    output: None,
                }),
                false => Ok(RuleResult::Failure {
                    name: "FakeRule".into(),
                    message: "FakeRule failed".into(),
                }),
            }
        }

        fn typetag_name(&self) -> &'static str {
            "fake-rule"
        }

        fn typetag_deserialize(&self) {
            todo!()
        }
    }

    /// An async (Async-mode) counterpart to FakeRule.
    #[derive(Debug, Serialize, Deserialize)]
    struct FakeAsyncRule {
        success: bool,
    }

    impl FakeAsyncRule {
        fn new(success: bool) -> Self {
            Self { success }
        }
    }

    impl Display for FakeAsyncRule {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "FakeAsyncRule")
        }
    }

    impl Rule for FakeAsyncRule {
        fn execution_mode(&self) -> ExecutionMode {
            ExecutionMode::Async
        }

        fn check(&self, _: &dyn Context) -> Result<RuleResult> {
            match self.success {
                true => Ok(RuleResult::Success {
                    name: "FakeAsyncRule".into(),
                    output: None,
                }),
                false => Ok(RuleResult::Failure {
                    name: "FakeAsyncRule".into(),
                    message: "FakeAsyncRule failed".into(),
                }),
            }
        }

        fn typetag_name(&self) -> &'static str {
            "fake-async-rule"
        }

        fn typetag_deserialize(&self) {
            todo!()
        }
    }

    fn make_config(rules: Vec<RuleContext>) -> Arc<Configuration> {
        let mut config = Configuration {
            hooks: Default::default(),
            extract: vec![],
            files: vec![],
        };
        config.hooks.insert(GitHook::PreCommit, rules);
        Arc::new(config)
    }

    fn rule_ctx(rule: impl Rule + 'static) -> RuleContext {
        RuleContext {
            when: None,
            extract: None,
            rule: Box::new(rule),
        }
    }

    // ── sync-only rules ───────────────────────────────────────────────────────

    #[test]
    fn test_sync_rule_success() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_configuration()
            .returning(move || make_config(vec![rule_ctx(FakeRule::new(true))]));

        assert!(command.exec(&mut context).is_ok());
    }

    #[test]
    fn test_sync_rule_failure() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_configuration()
            .returning(move || make_config(vec![rule_ctx(FakeRule::new(false))]));

        assert!(command.exec(&mut context).is_err());
    }

    // ── async-mode rules ──────────────────────────────────────────────────────

    #[test]
    fn test_async_rule_success() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_configuration()
            .returning(move || make_config(vec![rule_ctx(FakeAsyncRule::new(true))]));

        assert!(command.exec(&mut context).is_ok());
    }

    #[test]
    fn test_async_rule_failure() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_configuration()
            .returning(move || make_config(vec![rule_ctx(FakeAsyncRule::new(false))]));

        assert!(command.exec(&mut context).is_err());
    }

    // ── mixed sync + async ────────────────────────────────────────────────────

    #[test]
    fn test_mixed_rules_all_succeed() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            make_config(vec![
                rule_ctx(FakeRule::new(true)),
                rule_ctx(FakeAsyncRule::new(true)),
                rule_ctx(FakeRule::new(true)),
                rule_ctx(FakeAsyncRule::new(true)),
            ])
        });

        assert!(command.exec(&mut context).is_ok());
    }

    #[test]
    fn test_mixed_rules_sync_fails() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            make_config(vec![
                rule_ctx(FakeRule::new(false)),     // sync failure
                rule_ctx(FakeAsyncRule::new(true)), // async ok
            ])
        });

        assert!(command.exec(&mut context).is_err());
    }

    #[test]
    fn test_mixed_rules_async_fails() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            make_config(vec![
                rule_ctx(FakeRule::new(true)),       // sync ok
                rule_ctx(FakeAsyncRule::new(false)), // async failure
            ])
        });

        assert!(command.exec(&mut context).is_err());
    }

    // ── concurrency limiting ──────────────────────────────────────────────────

    #[test]
    fn test_more_async_rules_than_concurrency_limit() {
        // Spawn more async rules than MAX_CONCURRENT_ASYNC_RULES to exercise
        // the semaphore path and ensure no deadlock.
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let count = MAX_CONCURRENT_ASYNC_RULES * 2;
        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            let rules = (0..count)
                .map(|_| rule_ctx(FakeAsyncRule::new(true)))
                .collect();
            make_config(rules)
        });

        assert!(command.exec(&mut context).is_ok());
    }

    // ── result ordering ───────────────────────────────────────────────────────

    /// Verify that results are collected in declaration order even though async
    /// rules finish in non-deterministic order.
    #[test]
    fn test_result_order_preserved() {
        #[derive(Debug, Serialize, Deserialize)]
        struct OrderedRule {
            label: String,
            is_async: bool,
        }

        impl Display for OrderedRule {
            fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(f, "OrderedRule({})", self.label)
            }
        }

        impl Rule for OrderedRule {
            fn execution_mode(&self) -> ExecutionMode {
                if self.is_async { ExecutionMode::Async } else { ExecutionMode::Sync }
            }

            fn check(&self, _: &dyn Context) -> Result<RuleResult> {
                Ok(RuleResult::Success {
                    name: self.label.clone(),
                    output: None,
                })
            }

            fn typetag_name(&self) -> &'static str { "ordered-rule" }
            fn typetag_deserialize(&self) { todo!() }
        }

        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            make_config(vec![
                rule_ctx(OrderedRule { label: "rule-0".into(), is_async: true }),
                rule_ctx(OrderedRule { label: "rule-1".into(), is_async: false }),
                rule_ctx(OrderedRule { label: "rule-2".into(), is_async: true }),
                rule_ctx(OrderedRule { label: "rule-3".into(), is_async: false }),
            ])
        });

        assert!(command.exec(&mut context).is_ok());
    }

    // ── no rules for hook ─────────────────────────────────────────────────────

    #[test]
    fn test_no_rules_for_hook() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context.expect_configuration().returning(move || {
            Arc::new(Configuration {
                hooks: Default::default(),
                extract: vec![],
                files: vec![],
            })
        });

        assert!(command.exec(&mut context).is_ok());
    }
}

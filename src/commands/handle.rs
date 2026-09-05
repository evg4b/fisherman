use crate::commands::command::CliCommand;
use anyhow::{Result, anyhow};
use clap::Parser;
use fisherman_core::Context;
use fisherman_core::GitHook;
use fisherman_core::RuleResult;
use fisherman_core::hook_display;
use fisherman_core::{ExecutionMode, RuleExecutionPool};
use std::path::PathBuf;

#[derive(Debug, Parser)]
pub struct HandleCommand {
    /// The hook to handle
    #[arg(value_enum)]
    hook: GitHook,
    /// The commit message file path
    message: Option<String>,
}

impl CliCommand for HandleCommand {
    async fn exec(&self, context: &mut impl Context) -> Result<()> {
        if let Some(message) = &self.message {
            context.set_commit_msg_path(PathBuf::from(message));
        }

        let config = context.configuration();
        println!("{}", hook_display(&self.hook, config.files.clone()));

        let Some(rules) = config.hooks.get(&self.hook) else {
            println!("No rules found for hook {}", self.hook);
            return Ok(());
        };

        let (sync_indexes, async_indexes): (Vec<usize>, Vec<usize>) = (0..rules.len())
            .partition(|&index| rules[index].rule.execution_mode() == ExecutionMode::Sync);

        let pool = RuleExecutionPool::new();
        let mut handles = Vec::with_capacity(async_indexes.len() + 1);
        let hook = self.hook;

        // Sync rules are cheap, so they share a single task and one context,
        // running one after another while the async rules proceed alongside them.
        if !sync_indexes.is_empty() {
            let mut sync_context = context.extend(&[])?;
            let config = config.clone();

            handles.push(tokio::spawn(async move {
                let mut results = Vec::with_capacity(sync_indexes.len());
                for index in sync_indexes {
                    let rule = &config.hooks[&hook][index];
                    results.push((index, rule.check_rule(sync_context.as_mut()).await?));
                }

                anyhow::Ok(results)
            }));
        }

        for index in async_indexes {
            let mut rule_context = context.extend(&[])?;
            let config = config.clone();

            handles.push(pool.execute(async move {
                let rule = &config.hooks[&hook][index];

                anyhow::Ok(vec![(index, rule.check_rule(rule_context.as_mut()).await?)])
            }));
        }

        let mut results = Vec::<(usize, RuleResult)>::with_capacity(rules.len());
        for handle in handles {
            results.extend(handle.await??);
        }

        results.sort_by_key(|(index, _)| *index);

        for (_, rule) in &results {
            match rule {
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
                }
                RuleResult::Skipped { name } => {
                    println!("{name}: skipped");
                }
            }
        }

        if results
            .iter()
            .any(|(_, rule)| matches!(rule, RuleResult::Failure { .. }))
        {
            return Err(anyhow!("Hook failed"));
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use async_trait::async_trait;
    use fisherman_core::MockContext;
    use fisherman_core::Rule;
    use fisherman_core::{Configuration, RuleContext};
    use serde::{Deserialize, Serialize};
    use std::fmt::Display;
    use std::sync::{Arc, Mutex};

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

    #[async_trait]
    impl Rule for FakeRule {
        async fn check(&self, _: &dyn Context) -> Result<RuleResult> {
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
            todo!()
        }

        fn typetag_deserialize(&self) {
            todo!()
        }
    }

    /// Records which rules ran, in the order they ran.
    #[derive(Debug, Default, Serialize, Deserialize)]
    struct TrackedRule {
        name: String,
        success: bool,
        mode_is_async: bool,
        #[serde(skip)]
        executed: Arc<Mutex<Vec<String>>>,
    }

    impl Display for TrackedRule {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "TrackedRule({})", self.name)
        }
    }

    #[async_trait]
    impl Rule for TrackedRule {
        async fn check(&self, _: &dyn Context) -> Result<RuleResult> {
            // Give the scheduler a turn so interleaving is possible at all.
            tokio::task::yield_now().await;

            self.executed
                .lock()
                .expect("tracker poisoned")
                .push(self.name.clone());

            match self.success {
                true => Ok(RuleResult::Success {
                    name: self.name.clone(),
                    output: None,
                }),
                false => Ok(RuleResult::Failure {
                    name: self.name.clone(),
                    message: "failed".into(),
                }),
            }
        }

        fn execution_mode(&self) -> ExecutionMode {
            match self.mode_is_async {
                true => ExecutionMode::Async,
                false => ExecutionMode::Sync,
            }
        }

        fn typetag_name(&self) -> &'static str {
            todo!()
        }

        fn typetag_deserialize(&self) {
            todo!()
        }
    }

    /// `(name, succeeds, is_async)`
    type RuleSpec = (&'static str, bool, bool);

    fn tracked_context(
        specs: &'static [RuleSpec],
        executed: Arc<Mutex<Vec<String>>>,
    ) -> MockContext {
        let mut context = MockContext::new();
        context
            .expect_extend()
            .returning(|_| Ok(Box::new(MockContext::new()) as Box<dyn Context>));
        context.expect_configuration().returning(move || {
            let mut config = Configuration {
                hooks: Default::default(),
                extract: vec![],
                files: vec![],
            };

            config.hooks.insert(
                GitHook::PreCommit,
                specs
                    .iter()
                    .map(|(name, success, mode_is_async)| RuleContext {
                        when: None,
                        extract: None,
                        rule: Box::new(TrackedRule {
                            name: (*name).to_string(),
                            success: *success,
                            mode_is_async: *mode_is_async,
                            executed: executed.clone(),
                        }),
                    })
                    .collect(),
            );

            Arc::new(config)
        });

        context
    }

    fn handle_pre_commit() -> HandleCommand {
        HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        }
    }

    #[tokio::test]
    async fn runs_every_rule_in_both_scheduling_groups() {
        let executed = Arc::new(Mutex::new(Vec::new()));
        let mut context = tracked_context(
            &[
                ("sync-1", true, false),
                ("async-1", true, true),
                ("sync-2", true, false),
                ("async-2", true, true),
            ],
            executed.clone(),
        );

        assert!(handle_pre_commit().exec(&mut context).await.is_ok());

        let mut ran = executed.lock().unwrap().clone();
        ran.sort();
        assert_eq!(ran, vec!["async-1", "async-2", "sync-1", "sync-2"]);
    }

    #[tokio::test]
    async fn sync_rules_share_one_task_and_keep_their_order() {
        let executed = Arc::new(Mutex::new(Vec::new()));
        let mut context = tracked_context(
            &[
                ("sync-1", true, false),
                ("async-1", true, true),
                ("sync-2", true, false),
                ("sync-3", true, false),
            ],
            executed.clone(),
        );

        assert!(handle_pre_commit().exec(&mut context).await.is_ok());

        let ran = executed.lock().unwrap().clone();
        let sync_order: Vec<&String> = ran.iter().filter(|name| name.starts_with("sync")).collect();
        assert_eq!(sync_order, vec!["sync-1", "sync-2", "sync-3"]);
    }

    #[tokio::test]
    async fn failing_async_rule_fails_the_hook() {
        let executed = Arc::new(Mutex::new(Vec::new()));
        let mut context = tracked_context(
            &[("sync-1", true, false), ("async-1", false, true)],
            executed.clone(),
        );

        assert!(handle_pre_commit().exec(&mut context).await.is_err());
        assert_eq!(executed.lock().unwrap().len(), 2);
    }

    #[tokio::test]
    async fn test_run() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_extend()
            .returning(|_| Ok(Box::new(MockContext::new()) as Box<dyn Context>));
        context.expect_configuration().returning(move || {
            let mut config = Configuration {
                hooks: Default::default(),
                extract: vec![],
                files: vec![],
            };

            config.hooks.insert(
                GitHook::PreCommit,
                vec![RuleContext {
                    when: None,
                    extract: None,
                    rule: Box::new(FakeRule::new(true)),
                }],
            );

            Arc::new(config)
        });

        let result = command.exec(&mut context).await;

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_run_with_message() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
        context
            .expect_extend()
            .returning(|_| Ok(Box::new(MockContext::new()) as Box<dyn Context>));
        context.expect_configuration().returning(move || {
            let mut config = Configuration {
                hooks: Default::default(),
                extract: vec![],
                files: vec![],
            };

            config.hooks.insert(
                GitHook::PreCommit,
                vec![RuleContext {
                    when: None,
                    extract: None,
                    rule: Box::new(FakeRule::new(false)),
                }],
            );

            Arc::new(config)
        });

        let result = command.exec(&mut context).await;

        assert!(result.is_err());
    }
}

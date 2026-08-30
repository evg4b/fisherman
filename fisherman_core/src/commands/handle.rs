use crate::Context;
use crate::GitHook;
use crate::RuleResult;
use crate::commands::command::CliCommand;
use crate::rules::{ExecutionMode, MAX_CONCURRENT_ASYNC_RULES};
use crate::ui::hook_display;
use anyhow::{Result, anyhow};
use clap::Parser;
use std::path::PathBuf;
use std::sync::Arc;
use tokio::sync::Semaphore;

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

        let semaphore = Arc::new(Semaphore::new(MAX_CONCURRENT_ASYNC_RULES));
        let mut handles = Vec::with_capacity(rules.len());

        for index in 0..rules.len() {
            // Every task needs its own context, so each rule gets a detached copy.
            let mut rule_context = context.extend(&[])?;
            let config = config.clone();
            let semaphore = semaphore.clone();
            let hook = self.hook;

            handles.push(tokio::spawn(async move {
                let rule = &config.hooks[&hook][index];

                // Sync rules are cheap and stay unbounded; async rules do file I/O
                // or spawn processes, so their concurrency is capped.
                let _permit = match rule.rule.execution_mode() {
                    ExecutionMode::Async => Some(semaphore.acquire().await?),
                    ExecutionMode::Sync => None,
                };

                rule.check_rule(rule_context.as_mut()).await
            }));
        }

        let mut results = Vec::<RuleResult>::with_capacity(handles.len());
        for handle in handles {
            results.push(handle.await??);
        }

        for rule in &results {
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
            .any(|r| matches!(r, RuleResult::Failure { .. }))
        {
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
    use async_trait::async_trait;
    use crate::{Configuration, RuleContext};
    use serde::{Deserialize, Serialize};
    use std::fmt::Display;
    use std::sync::Arc;

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

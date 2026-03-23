use crate::commands::command::CliCommand;
use crate::ui::hook_display;
use crate::Context;
use crate::GitHook;
use crate::RuleResult;
use anyhow::{anyhow, Result};
use clap::Parser;
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
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        if let Some(message) = &self.message {
            context.set_commit_msg_path(PathBuf::from(message));
        }

        let config = context.configuration();
        println!("{}", hook_display(&self.hook, config.files.clone()));

        match config.hooks.get(&self.hook) {
            Some(rules) => {
                let results = rules
                    .iter()
                    .map(|r| r.check_rule(context))
                    .collect::<Result<Vec<RuleResult>>>()?;

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
            }
            None => println!("No rules found for hook {}", self.hook),
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
            todo!()
        }

        fn typetag_deserialize(&self) {
            todo!()
        }
    }

    #[test]
    fn test_run() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
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

        let result = command.exec(&mut context);

        assert!(result.is_ok());
    }

    #[test]
    fn test_run_with_message() {
        let command = HandleCommand {
            hook: GitHook::PreCommit,
            message: None,
        };

        let mut context = MockContext::new();
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

        let result = command.exec(&mut context);

        assert!(result.is_err());
    }
}

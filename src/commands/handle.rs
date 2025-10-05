use crate::commands::command::CliCommand;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::rules::{CompiledRule, Rule, RuleResult};
use crate::ui::hook_display;
use anyhow::Result;
use clap::Parser;
use rayon::prelude::*;
use std::path::PathBuf;
use std::process::exit;
use std::sync::{Arc, Mutex};

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

        let config = context.configuration()?;
        println!("{}", hook_display(&self.hook, config.files));

        match config.hooks.get(&self.hook) {
            Some(rules) => {
                let (sync_rules, async_rules) = compile_rules(context, rules)?;

                let mut results: Vec<RuleResult> = vec![];

                for rule in sync_rules {
                    results.push(rule.check(context)?);
                }

                let context_lock = Arc::new(Mutex::new(context));
                let async_results: Result<Vec<_>> = async_rules
                    .par_iter()
                    .map(|rule| {
                        let ctx = context_lock.lock().unwrap();
                        rule.check(&**ctx)
                    })
                    .collect();

                results.extend(async_results?);

                for rule in &results {
                    match rule {
                        RuleResult::Success { name, output } => {
                            println!("{name} executed successfully");
                            if let Some(value) = output && !value.is_empty() {
                                println!("{value}");
                            }
                        }
                        RuleResult::Failure { message, name } => {
                            eprintln!("{name}: {message}");
                        }
                    }
                }

                if results
                    .iter()
                    .any(|r| matches!(r, RuleResult::Failure { .. }))
                {
                    exit(1);
                }
            }
            None => println!("No rules found for hook {}", self.hook),
        }

        Ok(())
    }
}

type RulesBucket = Vec<Box<dyn CompiledRule>>;

fn compile_rules(context: &impl Context, rules: &[Rule]) -> Result<(RulesBucket, RulesBucket)> {
    let mut sync_rules: RulesBucket = vec![];
    let mut async_rules: RulesBucket = vec![];
    for rule in rules.iter() {
        if let Some(compiled_rule) = rule.compile(context)? {
            if compiled_rule.sync() {
                sync_rules.push(compiled_rule);
            } else {
                async_rules.push(compiled_rule);
            }
        }
    }
    Ok((sync_rules, async_rules))
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::rules::RuleParams;
    use std::collections::HashMap;

    #[test]
    fn test_compile_rules() -> Result<()> {
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let rules = vec![
            Rule {
                when: None,
                extract: None,
                params: RuleParams::CommitMessageRegex {
                    regex: "^Test".to_string(),
                },
            },
            Rule {
                when: None,
                extract: None,
                params: RuleParams::ShellScript {
                    env: None,
                    script: "echo 'Hello World'".to_string(),
                },
            },
        ];

        let (sync_rules, async_rules) = compile_rules(&context, &rules)?;

        assert_eq!(sync_rules.len(), 1);
        assert_eq!(async_rules.len(), 1);
        assert!(sync_rules[0].sync());
        assert!(!async_rules[0].sync());

        Ok(())
    }

}

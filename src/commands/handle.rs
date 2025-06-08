use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::rules::{CompiledRule, Rule, RuleResult};
use crate::ui::hook_display;
use anyhow::Result;
use std::process::exit;

type RulesBucket = Vec<Box<dyn CompiledRule>>;

pub fn handle_command(context: &impl Context, hook: &GitHook) -> Result<()> {
    let config = Configuration::load(context.repo_path())?;
    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            let (sync_rules, async_rules) = compile_rules(context, &config.extract, rules)?;

            let mut results: Vec<RuleResult> = vec![];

            for rule in sync_rules {
                results.push(rule.check(context)?);
            }

            for rule in async_rules {
                results.push(rule.check(context)?);
            }

            for rule in results.iter() {
                match rule {
                    RuleResult::Success { name, output } => {
                        println!("{} executed successfully", name);
                        if !output.is_empty() {
                            println!("{}", output);
                        }
                    }
                    RuleResult::Failure { message, name } => {
                        eprintln!("{}: {}", name, message);
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
        None => println!("No rules found for hook {}", hook),
    };

    Ok(())
}

fn compile_rules(
    context: &impl Context,
    extract: &Vec<String>,
    rules: &[Rule],
) -> Result<(RulesBucket, RulesBucket)> {
    let mut sync_rules: RulesBucket = vec![];
    let mut async_rules: RulesBucket = vec![];
    for rule in rules.iter() {
        if let Some(compiled_rule) = rule.compile(context, extract)? {
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

    #[test]
    fn test_compile_rules() -> Result<()> {
        let context = MockContext::new();
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

        let (sync_rules, async_rules) = compile_rules(&context, &vec![], &rules)?;

        assert_eq!(sync_rules.len(), 1);
        assert_eq!(async_rules.len(), 1);
        assert!(sync_rules[0].sync());
        assert!(!async_rules[0].sync());

        Ok(())
    }
}

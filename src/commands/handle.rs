use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::rules::RuleResult;
use crate::ui::hook_display;
use anyhow::Result;
use std::process::exit;

pub fn handle_command(context: &impl Context, hook: &GitHook) -> Result<()> {
    let config = Configuration::load(context.repo_path())?;
    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            let mut results: Vec<RuleResult> = vec![];
            for rule in rules.iter() {
                if let Some(compiled_rule) = rule.compile(context, &config.extract)? {
                    results.push(compiled_rule.check(context)?);
                }
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

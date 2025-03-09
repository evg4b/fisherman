use crate::common::BError;
use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::rules::{CompiledRule, RuleResult};
use crate::ui::hook_display;
use std::process::exit;

pub fn handle_command(context: &impl Context, hook: &GitHook) -> Result<(), BError> {
    let config = Configuration::load(context.repo_path())?;
    println!("{}", hook_display(hook, config.files));


    match config.hooks.get(hook) {
        Some(rules) => {
            let rules_to_exec = rules.iter()
                    .map(|rule| rule.compile(context, config.extract.clone()));


            let results: Vec<RuleResult> = rules_to_exec
                .map(|rule| rule.unwrap().check())
                .collect();

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

            if results.iter().any(|r| matches!(r, RuleResult::Failure { .. }))
            {
                exit(1);
            }
        }
        None => println!("No rules found for hook {}", hook),
    };

    Ok(())
}


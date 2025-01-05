use crate::common::BError;
use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::rules::{Rule, RuleResult};
use crate::ui::hook_display::hook_display;
use std::process::exit;

pub(crate) fn handle_command(context: &Context, hook: &GitHook) -> Result<(), BError> {
    let config = Configuration::load(context.repo_path())?;
    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            let rules_to_exec: Vec<Rule> =
                rules.iter().map(|rule| Rule::new(rule.clone())).collect();

            let results: Vec<RuleResult> = rules_to_exec
                .into_iter()
                .map(|rule| rule.exec())
                .collect();

            for rule in results.iter() {
                match rule {
                    RuleResult::Success { name } => {
                        println!("{} executed successfully", name);
                    }
                    RuleResult::Failure { message, name } => {
                        eprintln!("{}: {}", name, message);
                    }
                }
            }

            if results.iter().any(|r| matches!(r, RuleResult::Failure { .. })) {
                exit(1);
            }
        }
        None => println!("No rules found for hook {}", hook),
    };

    Ok(())
}

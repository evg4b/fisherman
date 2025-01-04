use crate::common::BError;
use crate::configuration::Configuration;
use crate::hooks::GitHook;
use crate::rules::{Rule, RuleResult};
use crate::ui::logo::hook_display;
use std::env;
use std::process::exit;

pub(crate) fn handle_command(hook: &GitHook) -> Result<(), BError> {
    let cwd = env::current_dir()?;

    let config = Configuration::load(&cwd)?;
    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            let rules_to_exec: Vec<Rule> =
                rules.iter().map(|rule| Rule::new(rule.clone())).collect();

            let results: Vec<RuleResult> = rules_to_exec
                .into_iter()
                // TODO: Handle errors
                .map(|rule| rule.exec())
                .collect();

            let failed: Vec<&RuleResult> = results.iter().filter(|rule| !rule.success).collect();

            if !failed.is_empty() {
                for rule in failed {
                    println!("Rule {} execution failed", rule.name);
                    println!("Output: {}", rule.message);
                }
                exit(1);
            }

            for rule in results {
                println!("Rule {} successfully executed", rule.name);
            }
        }
        None => {
            eprintln!("No rules found for hook {}", hook);
            return Ok(());
        }
    };

    Ok(())
}

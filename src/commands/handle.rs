use crate::common::BError;
use crate::configuration::Configuration;
use crate::hooks::GitHook;
use crate::rules::{Rule, RuleResult};
use std::env;
use std::process::exit;

pub(crate) fn handle_command(hook: &GitHook) -> Result<(), BError> {
    let current_dir = env::current_dir()?;

    let config = Configuration::load(&current_dir)?;
    println!("Configuration loaded from {:?}", config.files);

    match config.hooks.get(hook) {
        Some(rules) => {
            let rules_to_exec: Vec<Rule> = rules
                .iter()
                .map(|rule| Rule::new(rule.clone()))
                .collect();

            let results: Vec<RuleResult> = rules_to_exec
                .into_iter()
                // TODO: Handle errors
                .map(|rule| rule.exec())
                .collect();

            for rule in results {
                if rule.success {
                    println!("Rule {} successfully executed", rule.name);
                } else {
                    println!("Rule {} execution failed", rule.name);
                    println!("Output: {}", rule.message);
                    exit(1);
                }
            }
        }
        None => {
            eprintln!("No rules found for hook {}", hook);
            return Ok(());
        }
    };

    Ok(())
}

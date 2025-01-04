use crate::common::BError;
use crate::configuration::Configuration;
use crate::hooks::GitHook;
use std::env;

pub(crate) fn explain_command(hook: &GitHook) -> Result<(), BError> {
    let current_dir = env::current_dir()?;
    let config = Configuration::load(&current_dir)?;
    println!("Configuration loaded from {:?}", config.files);

    match config.hooks.get(hook) {
        Some(rules) => {
            rules.into_iter().for_each(|rule| {
                println!("{:?}", rule);
            });
        }
        None => {
            println!("No rules found for hook {}", hook);
        }
    };

    Ok(())
}

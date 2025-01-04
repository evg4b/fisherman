use crate::common::BError;
use crate::configuration::Configuration;
use crate::hooks::GitHook;
use std::env;
use crate::ui::logo::hook_display;

pub(crate) fn explain_command(hook: &GitHook) -> Result<(), BError> {
    let cwd = env::current_dir()?;
    let config = Configuration::load(&cwd)?;

    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            rules.iter().for_each(|rule| {
                println!("{:?}", rule);
            });
        }
        None => {
            println!("No rules found for hook {}", hook);
        }
    };

    Ok(())
}

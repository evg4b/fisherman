use crate::common::BError;
use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::hook_display;

pub fn explain_command(context: &impl Context, hook: &GitHook) -> Result<(), BError> {
    let config = Configuration::load(context.repo_path())?;

    println!("{}", hook_display(hook, config.files));

    match config.hooks.get(hook) {
        Some(rules) => {
            rules.iter().for_each(|rule| {
                println!("{}", rule);
            });
        }
        None => {
            println!("No rules found for hook {}", hook);
        }
    };

    Ok(())
}

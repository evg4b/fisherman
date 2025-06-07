use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::logo;
use anyhow::Result;
use clap::ValueEnum;
use std::fs;

pub fn install_command(
    context: &impl Context,
    hooks: Option<Vec<GitHook>>,
    force: bool,
) -> Result<()> {
    println!("{}", logo());

    let selected_hooks = match hooks {
        Some(hooks) => hooks,
        None => Configuration::load(context.repo_path())?
            .get_configured_hooks()
            .unwrap_or_else(|| GitHook::value_variants().into()),
    };

    fs::create_dir_all(context.hooks_dir())?;
    for hook in selected_hooks {
        hook.install(context, force)?;
        println!("Hook {} installed", hook);
    }

    Ok(())
}

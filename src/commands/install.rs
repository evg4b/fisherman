use crate::configuration::Configuration;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::logo;
use anyhow::Result;

pub fn install_command(context: &impl Context, hooks: Option<Vec<GitHook>>, force: bool) -> Result<()> {
    println!("{}", logo());

    let selected_hooks = match hooks {
        Some(hooks) => hooks,
        None => Configuration::load(context.repo_path())?
            .get_configured_hooks()
            .unwrap_or_else(GitHook::all),
    };

    for hook in selected_hooks {
        hook.install(context, force)?;
        println!("Hook {} installed", hook);
    }

    Ok(())
}

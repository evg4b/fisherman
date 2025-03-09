use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::logo;
use anyhow::Result;

pub fn init_command(context: &impl Context, hooks: Option<Vec<GitHook>>, force: bool) -> Result<()> {
    println!("{}", logo());

    for hook in hooks.unwrap_or(GitHook::all()) {
        hook.install(context, force)?;
        println!("Hook {} initialized", hook);
    }

    Ok(())
}

use crate::common::BError;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::logo;

pub fn init_command(context: &impl Context, hooks: Option<Vec<GitHook>>, force: bool) -> Result<(), BError> {
    println!("{}", logo());

    for hook in hooks.unwrap_or(GitHook::all()) {
        hook.install(context, force)?;
        println!("Hook {} initialized", hook);
    }

    Ok(())
}

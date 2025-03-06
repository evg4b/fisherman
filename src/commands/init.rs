use crate::common::BError;
use crate::context::Context;
use crate::hooks::files::{build_hook_content, override_hook, write_hook};
use crate::hooks::GitHook;
use crate::ui::logo::logo;

pub fn init_command(context: &impl Context, force: bool) -> Result<(), BError> {
    println!("{}", logo());

    for hook_name in GitHook::all() {
        if force {
            override_hook(context, hook_name, build_hook_content(context.bin(), hook_name))?;
        } else {
            write_hook(context, hook_name, build_hook_content(context.bin(), hook_name))?;
        }

        println!("Hook {} initialized", hook_name);
    }

    Ok(())
}

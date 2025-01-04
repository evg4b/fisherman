use crate::common::BError;
use crate::hooks::files::{build_hook_content, override_hook, write_hook};
use crate::hooks::GitHook;
use std::env;

pub fn init_command(force: bool) -> Result<(), BError> {
    let current_dir = env::current_dir()?;
    let bin = env::current_exe()?;

    for hook_name in GitHook::all() {
        if force {
            override_hook(&current_dir, hook_name, build_hook_content(&bin, hook_name))?;
        } else {
            write_hook(&current_dir, hook_name, build_hook_content(&bin, hook_name))?;
        }

        println!("Hook {} initialized", hook_name);
    }

    Ok(())
}

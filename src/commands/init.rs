use crate::common::BError;
use crate::hooks::files::{build_hook_content, override_hook, write_hook};
use crate::hooks::GitHook;
use crate::ui::logo::logo;
use std::env;

pub fn init_command(force: bool) -> Result<(), BError> {
    let (cwd, bin) = (env::current_dir()?, env::current_exe()?);

    println!("{}", logo());

    for hook_name in GitHook::all() {
        if force {
            override_hook(&cwd, hook_name, build_hook_content(&bin, hook_name))?;
        } else {
            write_hook(&cwd, hook_name, build_hook_content(&bin, hook_name))?;
        }

        println!("Hook {} initialized", hook_name);
    }

    Ok(())
}

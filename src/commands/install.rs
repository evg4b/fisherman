use crate::commands::command::CliCommand;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::logo;
use anyhow::Result;
use clap::{Parser, ValueEnum};
use std::fs;

#[derive(Debug, Parser)]
pub struct InstallCommand {
    /// List of hooks to install (if not provided, only the configured
    /// hooks will be installed or all hooks if no configuration is found)
    hooks: Option<Vec<GitHook>>,
    /// Force the initialization of the hooks (override existing hooks)
    #[arg(short, long)]
    force: bool,
}

impl CliCommand for InstallCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        println!("{}", logo());

        let selected_hooks = match self.hooks.as_ref() {
            Some(hooks) => hooks.clone(),
            None => context
                .configuration()?
                .get_configured_hooks()
                .unwrap_or_else(|| GitHook::value_variants().into()),
        };

        fs::create_dir_all(context.hooks_dir())?;
        for hook in selected_hooks {
            hook.install(context, self.force)?;
            println!("Hook {} installed", hook);
        }

        Ok(())
    }
}

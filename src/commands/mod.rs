mod explain;
mod handle;
mod install;

use std::path::PathBuf;
use crate::context::Context;
use crate::hooks::GitHook;
use anyhow::Result;
use clap::Subcommand;
pub use explain::explain_command;
pub use handle::handle_command;
pub use install::install_command;

#[derive(Subcommand)]
pub enum Commands {
    /// Install hooks for the repository
    Install {
        /// Force the initialization of the hooks (override existing hooks)
        #[arg(short, long)]
        force: bool,
        /// List of hooks to install (if not provided, only the configured
        /// hooks will be installed or all hooks if no configuration is found)
        hooks: Option<Vec<GitHook>>,
    },
    /// Handle a hook
    Handle {
        /// The hook to handle
        #[arg(value_enum)]
        hook: GitHook,
        /// The commit message file path
        message: Option<String>,
    },
    /// Explain a hook behavior
    Explain {
        /// The hook to explain
        #[arg(value_enum)]
        hook: GitHook,
    },
}

impl Commands {
    pub fn run(&self, context: &mut impl Context) -> Result<()> {
        match self {
            Commands::Install { force, hooks } => {
                install_command(context, hooks.clone(), *force)
            },
            Commands::Handle { hook, message } => {
                if let Some(message) = message {
                    context.set_commit_msg_path(PathBuf::from(message));
                }

                handle_command(context, hook)
            },
            Commands::Explain { hook } => {
                explain_command(context, hook)
            },
        }
    }
}

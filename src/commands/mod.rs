mod explain;
mod handle;
mod init;

use crate::context::Context;
use crate::hooks::GitHook;
use anyhow::Result;
use clap::Subcommand;
pub use explain::explain_command;
pub use handle::handle_command;
pub use init::init_command;

#[derive(Subcommand)]
pub enum Commands {
    /// Initialize hooks for the repository
    Init {
        /// Force the initialization of the hooks (override existing hooks)
        #[arg(short, long)]
        force: bool,
        hooks: Option<Vec<GitHook>>,
    },
    /// Handle a hook
    Handle {
        /// The hook to handle
        #[arg(value_enum)]
        hook: GitHook,
    },
    /// Explain a hook behavior
    Explain {
        /// The hook to explain
        #[arg(value_enum)]
        hook: GitHook,
    },
}

impl Commands {
    pub fn run(&self, context: &impl Context) -> Result<()> {
        match self {
            Commands::Init { force, hooks } => init_command(context, hooks.clone(), *force),
            Commands::Handle { hook } => handle_command(context, hook),
            Commands::Explain { hook } => explain_command(context, hook),
        }
    }
}

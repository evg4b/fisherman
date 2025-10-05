mod command;
mod explain;
mod handle;
mod install;

pub use crate::commands::command::CliCommand;
pub use crate::commands::explain::ExplainCommand;
pub use crate::commands::handle::HandleCommand;
pub use crate::commands::install::InstallCommand;
use crate::context::Context;
use anyhow::Result;
use clap::Subcommand;

#[derive(Subcommand, Debug)]
pub enum Command {
    /// Install hooks for the repository
    Install(InstallCommand),
    /// Handle a hook
    Handle(HandleCommand),
    /// Explain a hook behavior
    Explain(ExplainCommand),
}

impl CliCommand for Command {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        match self {
            Command::Install(cmd) => cmd.exec(context),
            Command::Handle(cmd) => cmd.exec(context),
            Command::Explain(cmd) => cmd.exec(context),
        }
    }
}

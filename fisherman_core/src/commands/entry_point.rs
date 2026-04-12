use crate::commands::{CliCommand, ExplainCommand, HandleCommand, InstallCommand};
use crate::ui::ABOUT;
use crate::Context;
use clap::{Parser, Subcommand};
use anyhow::Result;

#[derive(Parser, Debug)]
#[command(author, version, about = ABOUT, long_about=None)]
pub struct FishermanCli {
    #[clap(subcommand)]
    command: Command,
}

#[derive(Subcommand, Debug)]
pub enum Command {
    /// Install hooks for the repository
    Install(InstallCommand),
    /// Handle a hook
    Handle(HandleCommand),
    /// Explain a hook behavior
    Explain(ExplainCommand),
}

impl Default for FishermanCli {
    fn default() -> Self {
        FishermanCli::parse()
    }
}

impl FishermanCli {
    pub fn run(self, context: &mut impl Context) -> Result<()> {
        match &self.command {
            Command::Install(cmd) => cmd.exec(context),
            Command::Handle(cmd) => cmd.exec(context),
            Command::Explain(cmd) => cmd.exec(context),
        }
    }
}

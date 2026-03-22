use crate::commands::{CliCommand, ExplainCommand, HandleCommand, InstallCommand};
use crate::ui::ABOUT;
use crate::Context;
use clap::{Parser, Subcommand};

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

impl FishermanCli {
    pub fn exec(context: &mut impl Context) -> anyhow::Result<()> {
        let cli = FishermanCli::parse();
        match &cli.command {
            Command::Install(cmd) => cmd.exec(context),
            Command::Handle(cmd) => cmd.exec(context),
            Command::Explain(cmd) => cmd.exec(context),
        }
    }
}

mod commands;
mod ui;

use crate::commands::{CliCommand, Command};
use core::GitRepoContext;
use clap::Parser;
use std::env;
use ui::ABOUT;

#[derive(Parser, Debug)]
#[command(author, version, about = ABOUT, long_about=None)]
struct Cli {
    #[clap(subcommand)]
    command: Command,
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    let context = &mut GitRepoContext::new(env::current_dir()?)?;

    match cli.command.exec(context) {
        Ok(()) => Ok(()),
        Err(err) => {
            eprintln!("Error: {err}");
            std::process::exit(1);
        }
    }
}

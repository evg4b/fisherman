use crate::commands::{CliCommand, Command};
use crate::context::GitRepoContext;
use clap::{Parser, command};
use std::env;

mod commands;
mod configuration;
mod context;
mod hooks;
mod rules;
mod scripting;
mod templates;
mod ui;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about)]
struct Cli {
    #[clap(subcommand)]
    command: Command,
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    let context = &mut GitRepoContext::new(env::current_dir()?)?;

    match cli.command.exec(context) {
        Ok(_) => Ok(()),
        Err(err) => {
            eprintln!("Error: {}", err);
            std::process::exit(1);
        }
    }
}

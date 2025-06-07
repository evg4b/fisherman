use crate::commands::Commands;
use crate::context::GitRepoContext;
use clap::{command, Parser};
use std::env;

mod commands;
mod configuration;
mod context;
mod hooks;
mod rules;
mod scripting;
mod templates;
mod ui;

#[derive(Parser)]
#[command(author, version, about, long_about)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    let context = &mut GitRepoContext::new(env::current_dir()?)?;

    match cli.command.run(context) {
        Ok(_) => Ok(()),
        Err(err) => {
            eprintln!("Error: {}", err);
            std::process::exit(1);
        }
    }
}

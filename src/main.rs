use crate::commands::explain::explain_command;
use crate::commands::handle::handle_command;
use crate::commands::init::init_command;
use crate::common::BError;
use crate::context::Context;
use crate::hooks::GitHook;
use clap::{Parser, Subcommand};
use std::env;

mod commands;
mod common;
mod configuration;
mod hooks;
mod rules;
mod ui;
mod context;

#[derive(Parser)]
#[command(author, version, about, long_about)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize hooks for the repository
    Init {
        /// Force the initialization of the hooks (override existing hooks)
        #[arg(short, long)]
        force: bool,
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

fn main() -> Result<(), BError> {
    let cli = Cli::parse();

    let context = &Context::new(env::current_dir()?)?;
    println!("On branch: {}", context.current_branch()?);
    println!("Hooks dir: {:?}", context.hooks_dir());

    let result = match &cli.command {
        Commands::Init { force } => init_command(context, *force),
        Commands::Handle { hook } => handle_command(context, hook),
        Commands::Explain { hook } => explain_command(context, hook),
    };

    match result {
        Ok(_) => Ok(()),
        Err(err) => {
            eprintln!("Error: {}", err);
            std::process::exit(1);
        }
    }
}

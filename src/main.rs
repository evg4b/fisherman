use crate::commands::explain::explain_command;
use crate::commands::handle::handle_command;
use crate::commands::init::init_command;
use crate::common::BError;
use crate::hooks::GitHook;
use clap::{Parser, Subcommand};
use std::env;
use std::fmt::Display;

mod commands;
mod common;
mod configuration;
mod hooks;
mod rules;

#[derive(Parser)]
#[command(author, version, about, long_about)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

impl Display for Cli {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", logo())
    }
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

    match &cli.command {
        Commands::Init { force } => init_command(*force),
        Commands::Handle { hook } => handle_command(hook),
        Commands::Explain { hook } => explain_command(hook),
    }
}

fn logo() -> String {
    format!(
        r#"
 .d888  d8b          888
 d88P"  Y8P          888                        {:>30}
 888                 888
 888888 888 .d8888b  88888b.   .d88b.  888d888 88888b.d88b.   8888b.  88888b.
 888    888 88K      888 "88b d8P  Y8b 888P"   888 "888 "88b     "88b 888 "88b
 888    888 "Y8888b. 888  888 88888888 888     888  888  888 .d888888 888  888
 888    888      X88 888  888 Y8b.     888     888  888  888 888  888 888  888
 888    888  88888P' 888  888  "Y8888  888     888  888  888 "Y888888 888  888
"#,
        format!("Version: {}", env!("CARGO_PKG_VERSION"))
    )
}

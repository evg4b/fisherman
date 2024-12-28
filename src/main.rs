use crate::hooks::{build_hook_content, override_hook, write_hook, GitHook};
use clap::{Parser, Subcommand};
use std::env;
use std::error::Error;

pub mod hooks;

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

fn main() -> Result<(), Box<dyn Error>> {
    let cli = Cli::parse();

    match &cli.command {
        Commands::Init { force } => {
            let bin = env::current_exe().expect("Failed to get current executable path");
            let current_dir = env::current_dir().expect("Failed to get current working directory");
            for hook_name in GitHook::all() {
                if *force {
                    override_hook(&current_dir, hook_name, build_hook_content(&bin, hook_name))?;
                } else {
                    write_hook(&current_dir, hook_name, build_hook_content(&bin, hook_name))?;
                }
                println!("Hook {} initialized", hook_name);
            }

            Ok(())
        }
        Commands::Handle { hook } => {
            println!("Handling task {}", hook);
            Ok(())
        }
        Commands::Explain { hook } => {
            println!("Explain task {}", hook);
            Ok(())
        }
    }
}

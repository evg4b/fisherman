use std::env;
use clap::{Parser, Subcommand};
use std::error::Error;

pub mod hooks;
use crate::hooks::{backup_hook, build_hook_content, read_hooks, write_hook};

#[derive(Parser)]
#[command(author, version, about, long_about)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize hooks for the repository
    Init,
    /// Handle a hook task
    Handle,
}

fn main() -> Result<(), Box<dyn Error>> {
    let cli = Cli::parse();

    match &cli.command {
        Commands::Init => {
            let entries = read_hooks();
            let bin = env::current_exe().expect("Failed to get current executable path");

            for (hook_name, entry) in entries {
                if entry.exists() {
                    backup_hook(&entry)?;
                    println!("Backed up hook: {:?}", entry);
                }
                write_hook(&entry, build_hook_content(&bin, hook_name))?;
            }

            Ok(())
        }
        Commands::Handle => {
            println!("Handling task");
            Ok(())
        }
    }
}

use clap::{Parser, Subcommand};

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

fn main() {
    let cli = Cli::parse();

    match &cli.command {
        Commands::Init => {
            println!("Initializing project");
        }
        Commands::Handle => {
            println!("Handling task");
        }
    }
}
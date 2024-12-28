use clap::{Parser, Subcommand};
use std::{env, fs, io};
use std::error::Error;

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
            let cwd = env::current_dir()
                .expect("Failed to get current working directory");
            let bin_path = env::current_exe()
                .expect("Failed to get binary path");

            println!("Current working directory: {:?}", cwd);
            println!("Binary path: {:?}", bin_path);

            let hooks_directory = cwd.join(".git/hooks");
            println!("Hooks directory: {:?}", hooks_directory);
            if !hooks_directory.exists() {
                fs::create_dir(&hooks_directory)?;
            }

            let entries = fs::read_dir(hooks_directory)?
                .map(|res| res.map(|e| e.path()))
                .collect::<Result<Vec<_>, io::Error>>()?;

            for entry in entries {
                println!("Entry: {:?}", entry);
                let file = fs::read_to_string(&entry)?;
                println!("File:\n {:?}", file);
            }

            Ok(())
        }
        Commands::Handle => {
            println!("Handling task");
            Ok(())
        }
    }
}

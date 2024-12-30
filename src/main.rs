use crate::configuration::Configuration;
use crate::hooks::{build_hook_content, override_hook, write_hook, GitHook};
use clap::{Parser, Subcommand};
use std::env;
use std::error::Error;

mod configuration;
mod hooks;

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

const LOGO: &str = r#"
 .d888  d8b          888
 d88P"  Y8P          888                                 Version: {{FishermanVersion}}
 888                 888
 888888 888 .d8888b  88888b.   .d88b.  888d888 88888b.d88b.   8888b.  88888b.
 888    888 88K      888 "88b d8P  Y8b 888P"   888 "888 "88b     "88b 888 "88b
 888    888 "Y8888b. 888  888 88888888 888     888  888  888 .d888888 888  888
 888    888      X88 888  888 Y8b.     888     888  888  888 888  888 888  888
 888    888  88888P' 888  888  "Y8888  888     888  888  888 "Y888888 888  888
"#;

fn main() -> Result<(), Box<dyn Error>> {
    println!("{}", LOGO);

    let cli = Cli::parse();
    let current_dir = env::current_dir().expect("Failed to get current working directory");

    match &cli.command {
        Commands::Init { force } => {
            let bin = env::current_exe().expect("Failed to get current executable path");
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
            let config = Configuration::load(&current_dir)?;
            let item = config.hooks.unwrap_or_default();
            let rules = match item.get(hook) {
                Some(rules) => rules,
                None => {
                    eprintln!("No rules found for hook {}", hook);
                    return Ok(());
                }
            };

            println!("Handling task {} with rules {:?}", hook, rules);
            Ok(())
        }
        Commands::Explain { hook } => {
            let config = Configuration::load(&current_dir)?;
            let item = config.hooks.unwrap_or_default();
            let rules = match item.get(hook) {
                Some(rules) => rules,
                None => {
                    println!("No rules found for hook {}", hook);
                    return Ok(());
                }
            };

            println!("Handling task {} with rules {:?}", hook, rules);
            Ok(())
        }
    }
}

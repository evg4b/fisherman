mod commands;

use anyhow::Result;
use commands::FishermanCli;
use fisherman_core::GitRepoContext;
use std::env;

#[tokio::main]
async fn main() -> Result<()> {
    let cli = FishermanCli::default();
    let cwd = env::current_dir()?;
    let mut context = GitRepoContext::new(cwd)?;

    if let Err(err) = cli.run(&mut context).await {
        eprintln!("Error: {err}");
        std::process::exit(1);
    }

    Ok(())
}

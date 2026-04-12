use anyhow::Result;
use fisherman_core::FishermanCli;
use fisherman_core::GitRepoContext;
use std::env;

fn main() -> Result<()> {
    let cli = FishermanCli::default();
    let mut context = GitRepoContext::new(env::current_dir()?)?;
    match cli.run(&mut context) {
        Ok(()) => Ok(()),
        Err(err) => {
            eprintln!("Error: {err}");
            std::process::exit(1);
        }
    }
}

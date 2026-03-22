use anyhow::Result;
use fisherman_core::FishermanCli;
use fisherman_core::GitRepoContext;
use std::env;

fn main() -> Result<()> {
    let mut context = GitRepoContext::new(env::current_dir()?)?;
    match FishermanCli::exec(&mut context) {
        Ok(()) => Ok(()),
        Err(err) => {
            eprintln!("Error: {err}");
            std::process::exit(1);
        }
    }
}

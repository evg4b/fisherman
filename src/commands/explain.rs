use crate::commands::command::CliCommand;
use crate::context::Context;
use crate::hooks::GitHook;
use crate::ui::hook_display;
use anyhow::Result;
use clap::Parser;

#[derive(Debug, Parser)]
pub struct ExplainCommand {
    /// The hook to explain
    #[arg(value_enum)]
    hook: GitHook,
}

impl CliCommand for ExplainCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        let config = context.configuration()?;

        println!("{}", hook_display(&self.hook, config.files));

        match config.hooks.get(&self.hook) {
            Some(rules) => {
                for rule in rules {
                    println!("{rule}");
                }
            }
            None => println!("No rules found for hook {}", self.hook),
        }

        Ok(())
    }
}

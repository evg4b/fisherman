use crate::commands::command::CliCommand;
use crate::ui::hook_display;
use anyhow::Result;
use clap::Parser;
use core::context::Context;
use core::hooks::GitHook;
use core::rules::rule::RuleResult;
use std::path::PathBuf;
use std::process::exit;

#[derive(Debug, Parser)]
pub struct HandleCommand {
    /// The hook to handle
    #[arg(value_enum)]
    hook: GitHook,
    /// The commit message file path
    message: Option<String>,
}

impl CliCommand for HandleCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        if let Some(message) = &self.message {
            context.set_commit_msg_path(PathBuf::from(message));
        }

        let config = context.configuration()?;
        println!("{}", hook_display(&self.hook, config.files));

        match config.hooks.get(&self.hook) {
            Some(rules) => {
                let results = rules.iter()
                    .map(|r| r.check(context))
                    .collect::<Result<Vec<RuleResult>>>()?;

                for rule in &results {
                    match rule {
                        RuleResult::Success { name, output } => {
                            println!("{name} executed successfully");
                            if let Some(value) = output && !value.is_empty() {
                                println!("{value}");
                            }
                        }
                        RuleResult::Failure { message, name } => {
                            eprintln!("{name}: {message}");
                        },
                        RuleResult::Skipped => {
                            println!("skipped");
                        }
                    }
                }

                if results
                    .iter()
                    .any(|r| matches!(r, RuleResult::Failure { .. }))
                {
                    exit(1);
                }
            }
            None => println!("No rules found for hook {}", self.hook),
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
}

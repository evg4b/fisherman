use crate::commands::command::CliCommand;
use core::Context;
use core::GitHook;
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
                    println!("{}", rule);
                }
            }
            None => println!("No rules found for hook {}", self.hook),
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use core::Configuration;
    use core::MockContext;
    use core::CommitMessageRegexRule;
    use core::Rule;
    use core::t;
    use std::collections::HashMap;

    #[test]
    fn test_exec_no_rules_for_hook() -> Result<()> {
        let cmd = ExplainCommand {
            hook: GitHook::PreCommit,
        };

        let config = Configuration {
            hooks: HashMap::new(),
            extract: vec![],
            files: vec![],
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration().return_once(move || Ok(config));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());

        Ok(())
    }

    #[test]
    fn test_exec_with_rules_for_hook() -> Result<()> {
        let cmd = ExplainCommand {
            hook: GitHook::PreCommit,
        };

        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^feat"),
        };

        let config = Configuration {
            hooks: HashMap::from([(GitHook::PreCommit, vec![
                Box::new(rule) as Box<dyn Rule>,
            ])]),
            extract: vec![],
            files: vec![],
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration().return_once(move || Ok(config));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());

        Ok(())
    }

    #[test]
    fn test_exec_configuration_error() {
        let cmd = ExplainCommand {
            hook: GitHook::PreCommit,
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration()
            .return_once(|| Err(anyhow::anyhow!("Config error")));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_err());
    }
}

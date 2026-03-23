use crate::commands::command::CliCommand;
use crate::ui::hook_display;
use anyhow::Result;
use clap::Parser;
use crate::Context;
use crate::GitHook;

#[derive(Debug, Parser)]
pub struct ExplainCommand {
    /// The hook to explain
    #[arg(value_enum)]
    hook: GitHook,
}

impl CliCommand for ExplainCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        let config = context.configuration();

        println!("{}", hook_display(&self.hook, config.files.clone()));

        match config.hooks.get(&self.hook) {
            Some(rules) => {
                for rule in rules {
                    println!("{}", rule.rule);
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
    use crate::{t, RuleContext};
    use crate::CommitMessageRegexRule;
    use crate::Configuration;
    use crate::MockContext;
    use crate::Rule;
    use std::collections::HashMap;
    use std::sync::Arc;

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
        ctx.expect_configuration().return_once(move || Arc::new(config));

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
                RuleContext {
                    extract: None,
                    when: None,
                    rule: Box::new(rule) as Box<dyn Rule>
                },
            ])]),
            extract: vec![],
            files: vec![],
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration().return_once(move || Arc::new(config));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());

        Ok(())
    }
}

use crate::commands::command::CliCommand;
use core::context::Context;
use core::hooks::GitHook;
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
                    let json = serde_json::to_value(rule.as_ref())
                        .unwrap_or(serde_json::Value::Null);
                    let type_name = json.get("type")
                        .and_then(|v| v.as_str())
                        .unwrap_or("unknown")
                        .replace('-', " ");
                    println!("{}: {}", type_name, json);
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
    use core::configuration::Configuration;
    use core::context::MockContext;
    use core::rules::commit_message_regex::CommitMessageRegexRule;
    use core::rules::rule::Rule;
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

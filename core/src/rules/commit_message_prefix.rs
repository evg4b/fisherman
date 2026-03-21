use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use rules_derive::ConditionalRule as ConditionalRuleDerive;

static MESSAGE_PREFIX_RULE_NAME: &str = "message-prefix";

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct CommitMessagePrefixRule {
    pub when: Option<Expression>,
    pub prefix: TemplateString,
}

impl std::fmt::Display for CommitMessagePrefixRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Commit message must start with: {}", self.prefix)
    }
}

#[typetag::serde(name = "message-prefix")]
impl Rule for CommitMessagePrefixRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: MESSAGE_PREFIX_RULE_NAME.to_string(),
            });
        }

        let prefix = compile_tmpl(ctx, &self.prefix, &[])?;
        let commit_msg = ctx.commit_msg()?;

        match commit_msg.starts_with(&prefix) {
            true => Ok(RuleResult::Success {
                name: MESSAGE_PREFIX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: MESSAGE_PREFIX_RULE_NAME.to_string(),
                message: format!("Commit message must start with: {}", prefix),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use assert2::assert;
    use std::collections::HashMap;

    #[test]
    fn test_commit_message_prefix_success() -> anyhow::Result<()> {
        let rule = CommitMessagePrefixRule {
            when: None,
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "message-prefix");
        assert_eq!(output, None);

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_failure() -> anyhow::Result<()> {
        let rule = CommitMessagePrefixRule {
            when: None,
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("fix: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "message-prefix");
        assert_eq!(message, "Commit message must start with: feat");

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_variables_error() -> anyhow::Result<()> {
        let rule = CommitMessagePrefixRule {
            when: None,
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: message".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_commit_msg_error() -> anyhow::Result<()> {
        let rule = CommitMessagePrefixRule {
            when: None,
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Err(anyhow::anyhow!("Commit message error")));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_display() {
        let rule = CommitMessagePrefixRule { when: None, prefix: "feat:".into() };
        assert_eq!(format!("{}", rule), "Commit message must start with: `feat:`");
    }
}

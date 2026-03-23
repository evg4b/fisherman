use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{Rule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

static MESSAGE_PREFIX_RULE_NAME: &str = "message-prefix";

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct CommitMessagePrefixRule {
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
    use anyhow::anyhow;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = CommitMessagePrefixRule {
            prefix: t!("feat:"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"prefix":"feat:"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: CommitMessagePrefixRule = serde_json::from_str(r#"{"prefix":"feat:"}"#)?;

        assert_eq!(config.prefix, t!("feat:"));

        Ok(())
    }


    #[test]
    fn test_commit_message_prefix_success() -> Result<()> {
        let rule = CommitMessagePrefixRule {
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "message-prefix");
        assert_eq!(output, None);

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_failure() -> Result<()> {
        let rule = CommitMessagePrefixRule {
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("fix: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "message-prefix");
        assert_eq!(message, "Commit message must start with: feat");

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_variables_error() -> Result<()> {
        let rule = CommitMessagePrefixRule {
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: message".to_string()));
        ctx.expect_variables()
            .returning(|| Err(anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_commit_message_prefix_commit_msg_error() -> Result<()> {
        let rule = CommitMessagePrefixRule {
            prefix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Err(anyhow!("Commit message error")));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_display() {
        let rule = CommitMessagePrefixRule {
            prefix: "feat:".into(),
        };
        assert_eq!(
            format!("{}", rule),
            "Commit message must start with: `feat:`"
        );
    }
}

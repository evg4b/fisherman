use crate::context::Context;
use crate::rules::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use rules_derive::ConditionalRule as ConditionalRuleDerive;

static MESSAGE_SUFFIX_RULE_NAME: &str = "message-suffix";

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct CommitMessageSuffixRule {
    pub when: Option<Expression>,
    pub suffix: TemplateString,
}

impl std::fmt::Display for CommitMessageSuffixRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Commit message must end with: {}", self.suffix)
    }
}

#[typetag::serde(name = "message-suffix")]
impl Rule for CommitMessageSuffixRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: MESSAGE_SUFFIX_RULE_NAME.to_string(),
            });
        }

        let suffix = self.suffix.compile(&ctx.variables(&[])?)?;
        let commit_msg = ctx.commit_msg()?;

        match commit_msg.ends_with(&suffix) {
            true => Ok(RuleResult::Success {
                name: MESSAGE_SUFFIX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: MESSAGE_SUFFIX_RULE_NAME.to_string(),
                message: format!("Commit message must end with: {}", suffix),
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
    fn test_commit_message_suffix() {
        let rule = CommitMessageSuffixRule {
            when: None,
            suffix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("my commit message feat".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx).unwrap();
        match result {
            RuleResult::Success { name, .. } => {
                assert_eq!(name, "message-suffix");
            }
            _ => panic!("Expected Success"),
        }
    }

    #[test]
    fn test_commit_message_suffix_failure() {
        let rule = CommitMessageSuffixRule {
            when: None,
            suffix: t!("feat"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("my commit message".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx).unwrap();
        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "message-suffix");
                assert_eq!(message, "Commit message must end with: feat");
            }
            _ => panic!("Expected Failure"),
        }
    }

    #[test]
    fn test_commit_message_suffix_variables_error() {
        let rule = CommitMessageSuffixRule {
            when: None,
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("message suffix".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_commit_message_suffix_commit_msg_error() {
        let rule = CommitMessageSuffixRule {
            when: None,
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Err(anyhow::anyhow!("Commit message error")));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = CommitMessageSuffixRule {
            when: None,
            suffix: " [skip-ci]".into(),
        };
        assert_eq!(
            format!("{}", rule),
            "Commit message must end with: ` [skip-ci]`"
        );
    }
}

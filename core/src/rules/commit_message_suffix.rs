use crate::context::Context;
use crate::rules::rule::{Rule, RuleResult};
use crate::rules::{CompiledRule, RuleResultOld};
use crate::templates::TemplateString;
use anyhow::Result;

static MESSAGE_SUFFIX_RULE_NAME: &str = "message-suffix";

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct CommitMessageSuffixRule {
    pub suffix: TemplateString,
}

#[typetag::serde(name = "message-suffix")]
impl Rule for CommitMessageSuffixRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
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

#[derive(Debug)]
pub struct CommitMessageSuffix {
    name: String,
    suffix: TemplateString,
}

impl CommitMessageSuffix {
    pub fn new(name: String, suffix: TemplateString) -> Self {
        Self { name, suffix }
    }
}

impl CompiledRule for CommitMessageSuffix {
    fn is_sequential(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResultOld> {
        let suffix = self.suffix.compile(&ctx.variables(&[])?)?;
        let commit_msg = ctx.commit_msg()?;

        match commit_msg.ends_with(&suffix) {
            true => Ok(RuleResultOld::Success {
                name: self.name.clone(),
                output: None,
            }),
            false => Ok(RuleResultOld::Failure {
                name: self.name.clone(),
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
        let rule = CommitMessageSuffixRule { suffix: t!("feat") };
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
        let rule = CommitMessageSuffixRule { suffix: t!("feat") };
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
}

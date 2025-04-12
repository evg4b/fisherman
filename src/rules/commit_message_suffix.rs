use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

#[derive(Debug)]
pub struct CommitMessageSuffix {
    name: String,
    prefix: TemplateString,
}

impl CommitMessageSuffix {
    pub fn new(name: String, prefix: TemplateString) -> Self {
        Self {
            name,
            prefix,
        }
    }
}

impl CompiledRule for CommitMessageSuffix {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let processed_prefix = self.prefix.to_string()?;
        let commit_message = ctx.commit_msg()?;
        if commit_message.ends_with(&processed_prefix) {
            Ok(RuleResult::Success {
                name: self.name.clone(),
                output: processed_prefix,
            })
        } else {
            Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Commit message does not start with prefix: {}", processed_prefix),
            })
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;
    use assertor::{assert_that, EqualityAssertion};

    #[test]
    fn test_commit_message_suffix() {
        let rule = CommitMessageSuffix::new(
            "commit_message_suffix".to_string(),
            tmpl!("feat"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg().returning(|| Ok("my commit message feat".to_string()));
        
        let RuleResult::Success{ name, .. } = rule.check(&ctx).unwrap() else { panic!() };

        assert_that!(name).is_equal_to("commit_message_suffix".to_string());
    }
    
    #[test]
    fn test_commit_message_suffix_failure() {
        let rule = CommitMessageSuffix::new(
            "commit_message_suffix".to_string(),
            tmpl!("feat"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg().returning(|| Ok("my commit message".to_string()));
        
        let RuleResult::Failure{ name, message } = rule.check(&ctx).unwrap() else { panic!() };

        assert_that!(name).is_equal_to("commit_message_suffix".to_string());
        assert_that!(message).is_equal_to("Commit message does not start with prefix: feat".to_string());
    }
}
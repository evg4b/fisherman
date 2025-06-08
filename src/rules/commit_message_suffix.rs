use crate::context::Context;
use crate::rules::helpers::check_suffix;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

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
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        match check_suffix(&self.suffix, &ctx.commit_msg()?)? {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: self.suffix.to_string()?,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Commit message does not end with suffix: {}",
                    self.suffix.to_string()?
                ),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;
    use assertor::{EqualityAssertion, assert_that};

    #[test]
    fn test_commit_message_suffix() {
        let rule = CommitMessageSuffix::new("commit_message_suffix".to_string(), tmpl!("feat"));
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("my commit message feat".to_string()));

        let RuleResult::Success { name, .. } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_suffix".to_string());
    }

    #[test]
    fn test_commit_message_suffix_failure() {
        let rule = CommitMessageSuffix::new("commit_message_suffix".to_string(), tmpl!("feat"));
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("my commit message".to_string()));

        let RuleResult::Failure { name, message } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_suffix".to_string());
        assert_that!(message)
            .is_equal_to("Commit message does not end with suffix: feat".to_string());
    }
}

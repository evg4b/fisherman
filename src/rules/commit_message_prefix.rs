use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateStringLegacy;
use anyhow::Result;
use crate::rules::helpers::check_prefix;

#[derive(Debug)]
pub struct CommitMessagePrefix {
    name: String,
    prefix: TemplateStringLegacy,
}

impl CommitMessagePrefix {
    pub fn new(name: String, prefix: TemplateStringLegacy) -> Self {
        Self { name, prefix }
    }
}

impl CompiledRule for CommitMessagePrefix {
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        match check_prefix(&self.prefix, &ctx.commit_msg()?)? {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: self.prefix.to_string()?,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Commit message does not start with prefix: {}",
                    self.prefix.to_string()?
                ),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl_legacy;
    use assertor::{assert_that, EqualityAssertion};

    #[test]
    fn test_commit_message_prefix() {
        let rule = CommitMessagePrefix::new(
            "commit_message_prefix".to_string(),
            tmpl_legacy!("feat"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: my commit message".to_string()));

        let RuleResult::Success { name, output } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_prefix".to_string());
        assert_that!(output).is_equal_to("feat".to_string());
    }

    #[test]
    fn test_commit_message_prefix_failure() {
        let rule = CommitMessagePrefix::new(
            "commit_message_prefix".to_string(),
            tmpl_legacy!("feat".to_string()),
        );
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("fix: my commit message".to_string()));

        let RuleResult::Failure { name, message } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_prefix".to_string());
        assert_that!(message)
            .is_equal_to("Commit message does not start with prefix: feat".to_string());
    }
    
    #[test]
    fn test_sync() {
        let rule = CommitMessagePrefix::new("commit_message_prefix".to_string(), tmpl_legacy!("feat"));
        assert!(rule.sync());
    }
}

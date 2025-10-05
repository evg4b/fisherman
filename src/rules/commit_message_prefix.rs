use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

#[derive(Debug)]
pub struct CommitMessagePrefix {
    name: String,
    prefix: TemplateString,
}

impl CommitMessagePrefix {
    pub fn new(name: String, prefix: TemplateString) -> Self {
        Self { name, prefix }
    }
}

impl CompiledRule for CommitMessagePrefix {
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let prefix = compile_tmpl(ctx, &self.prefix, &[])?;
        let commit_msg = ctx.commit_msg()?;

        match commit_msg.starts_with(&prefix) {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Commit message does not start with prefix: {}", prefix),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use assertor::{EqualityAssertion, assert_that};
    use std::collections::HashMap;

    #[test]
    fn test_commit_message_prefix() {
        let rule = CommitMessagePrefix::new("commit_message_prefix".to_string(), t!("feat"));
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("feat: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Success { name, output } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_prefix".to_string());
        assert_eq!(output, None);
    }

    #[test]
    fn test_commit_message_prefix_failure() {
        let rule = CommitMessagePrefix::new("commit_message_prefix".to_string(), t!("feat"));
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg()
            .returning(|| Ok("fix: my commit message".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Failure { name, message } = rule.check(&ctx).unwrap() else {
            panic!()
        };

        assert_that!(name).is_equal_to("commit_message_prefix".to_string());
        assert_that!(message)
            .is_equal_to("Commit message does not start with prefix: feat".to_string());
    }

    #[test]
    fn test_sync() {
        let rule = CommitMessagePrefix::new("commit_message_prefix".to_string(), t!("feat"));
        assert!(rule.sync());
    }
}

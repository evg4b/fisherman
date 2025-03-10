use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use std::collections::HashMap;
use crate::templates::replace_in_string;
use anyhow::Result;

#[derive(Debug)]
pub struct CommitMessageSuffix {
    name: String,
    prefix: String,
    variables: HashMap<String, String>,
}

impl CommitMessageSuffix {
    pub fn new(name: String, prefix: String, variables: HashMap<String, String>) -> Self {
        Self {
            name,
            prefix,
            variables,
        }
    }
}

impl CompiledRule for CommitMessageSuffix {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let processed_prefix = replace_in_string(&self.prefix, &self.variables)?;
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
    use assertor::{assert_that, EqualityAssertion};
    use super::*;
    use crate::context::MockContext;

    #[test]
    fn test_commit_message_suffix() {
        let rule = CommitMessageSuffix::new(
            "commit_message_suffix".to_string(),
            "feat".to_string(),
            HashMap::new(),
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
            "feat".to_string(),
            HashMap::new(),
        );
        let mut ctx = MockContext::new();
        ctx.expect_commit_msg().returning(|| Ok("my commit message".to_string()));
        
        let RuleResult::Failure{ name, message } = rule.check(&ctx).unwrap() else { panic!() };

        assert_that!(name).is_equal_to("commit_message_suffix".to_string());
        assert_that!(message).is_equal_to("Commit message does not start with prefix: feat".to_string());
    }
}
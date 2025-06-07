use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;

pub struct BranchNamePrefix {
    name: String,
    prefix: TemplateString,
}

impl BranchNamePrefix {
    pub fn new(name: String, prefix: TemplateString) -> Self {
        Self { name, prefix }
    }
}

impl CompiledRule for BranchNamePrefix {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let processed_prefix = self.prefix.to_string()?;
        let branch_name = ctx.current_branch()?;
        if branch_name.starts_with(&processed_prefix) {
            Ok(RuleResult::Success {
                name: self.name.clone(),
                output: processed_prefix,
            })
        } else {
            Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Branch name does not start with prefix: {}",
                    processed_prefix
                ),
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
    fn test_branch_name_prefix() -> anyhow::Result<()> {
        let rule = BranchNamePrefix::new(
            "branch_name_prefix".to_string(),
            tmpl!("feat/"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("feat/my-feature".to_string()));

        let RuleResult::Success { name, .. } = rule.check(&ctx)? else { panic!() };

        assert_that!(name).is_equal_to("branch_name_prefix".to_string());

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_failure() -> anyhow::Result<()> {
        let rule = BranchNamePrefix::new(
            "branch_name_prefix".to_string(),
            tmpl!("feat/"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("bugfix/my-feature".to_string()));

        let RuleResult::Failure { name, message } = rule.check(&ctx)? else { panic!() };

        assert_that!(name).is_equal_to("branch_name_prefix".to_string());
        assert_that!(message).is_equal_to("Branch name does not start with prefix: feat/".to_string());

        Ok(())
    }
}
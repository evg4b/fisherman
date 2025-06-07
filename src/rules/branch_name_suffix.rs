use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;

pub struct BranchNameSuffix {
    name: String,
    suffix: TemplateString,
}

impl BranchNameSuffix {
    pub fn new(name: String, suffix: TemplateString) -> Self {
        Self { name, suffix }
    }
}

impl CompiledRule for BranchNameSuffix {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let processed_prefix = self.suffix.to_string()?;
        let branch_name = ctx.current_branch()?;
        if branch_name.ends_with(&processed_prefix) {
            Ok(RuleResult::Success {
                name: self.name.clone(),
                output: processed_prefix,
            })
        } else {
            Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Branch name does not end with suffix: {}", processed_prefix),
            })
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;

    #[test]
    fn test_branch_name_suffix_success() -> anyhow::Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));

        let result =
            BranchNameSuffix::new("Test Rule".to_string(), tmpl!("feature")).check(&ctx)?;

        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn test_branch_name_suffix_failure() -> anyhow::Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));

        let result = BranchNameSuffix::new("Test Rule".to_string(), tmpl!("suffix")).check(&ctx)?;

        assert!(matches!(result, RuleResult::Failure { .. }));

        Ok(())
    }
}

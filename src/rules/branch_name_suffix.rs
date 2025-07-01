use crate::context::Context;
use crate::rules::helpers::check_suffix;
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
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        match check_suffix(ctx, &self.suffix, &ctx.current_branch()?)? {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: self.suffix.to_string(&ctx.variables(&[])?)?,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Branch name does not end with suffix: {}",
                    self.suffix.to_string(&ctx.variables(&[])?)?
                ),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;
    use super::*;
    use crate::context::MockContext;
    use crate::t;

    #[test]
    fn test_branch_name_suffix_success() -> anyhow::Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result =
            BranchNameSuffix::new("Test Rule".to_string(), t!("feature")).check(&ctx)?;

        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn test_branch_name_suffix_failure() -> anyhow::Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result =
            BranchNameSuffix::new("Test Rule".to_string(), t!("suffix")).check(&ctx)?;

        assert!(matches!(result, RuleResult::Failure { .. }));

        Ok(())
    }

    #[test]
    fn test_sync() {
        let rule = BranchNameSuffix::new("Test Rule".to_string(), t!("suffix"));
        assert!(rule.sync());
    }
}

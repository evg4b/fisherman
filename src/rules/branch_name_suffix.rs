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
        match check_suffix(&self.suffix, &ctx.current_branch()?)? {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: self.suffix.to_string()?,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "Branch name does not end with suffix: {}",
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

use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
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
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let prefix = compile_tmpl(ctx, &self.prefix, &[])?;
        let branch_name = ctx.current_branch()?;

        match branch_name.starts_with(&prefix) {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Branch name must start with: {}", prefix),
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
    fn test_branch_name_prefix() -> anyhow::Result<()> {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Success { name, .. } = rule.check(&ctx)? else {
            panic!()
        };

        assert!(name == "branch_name_prefix");

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_failure() -> anyhow::Result<()> {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Failure { name, message } = rule.check(&ctx)? else {
            panic!()
        };

        assert!(name == "branch_name_prefix");
        assert!(message == "Branch name must start with: feat/");

        Ok(())
    }

    #[test]
    fn test_sync() {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));

        assert!(rule.sync());
    }
}

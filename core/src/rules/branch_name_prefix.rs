use serde::Deserialize;
use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{CompiledRule, RuleResultOld};
use crate::rules::rule::Rule;
use crate::templates::TemplateString;


#[derive(Debug, Deserialize, serde::Serialize)]
struct BranchNamePrefixRule {
    pub prefix: TemplateString,
}

#[typetag::serde(name = "branch-name-prefix")]
impl Rule for BranchNamePrefixRule {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<crate::rules::rule::RuleResult> {
        todo!()
    }
}


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
    fn is_sequential(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResultOld> {
        let prefix = compile_tmpl(ctx, &self.prefix, &[])?;
        let branch_name = ctx.current_branch()?;

        match branch_name.starts_with(&prefix) {
            true => Ok(RuleResultOld::Success {
                name: self.name.clone(),
                output: None,
            }),
            false => Ok(RuleResultOld::Failure {
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

        let result = rule.check(&ctx)?;
        let RuleResultOld::Success { name, .. } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "branch_name_prefix");

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

        let result = rule.check(&ctx)?;
        let RuleResultOld::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "branch_name_prefix");
        assert_eq!(message, "Branch name must start with: feat/");

        Ok(())
    }

    #[test]
    fn test_is_sequential() {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));

        assert!(rule.is_sequential());
    }

    #[test]
    fn test_branch_name_prefix_variables_error() {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_branch_name_prefix_branch_error() {
        let rule = BranchNamePrefix::new("branch_name_prefix".to_string(), t!("feat/"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Err(anyhow::anyhow!("Branch error")));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }
}

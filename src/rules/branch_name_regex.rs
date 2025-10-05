use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use regex::Regex;

pub struct BranchNameRegex {
    name: String,
    expression: TemplateString,
}

impl BranchNameRegex {
    pub fn new(name: String, expression: TemplateString) -> Self {
        Self { name, expression }
    }
}

impl CompiledRule for BranchNameRegex {
    fn sync(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let expression = Regex::new(&compile_tmpl(ctx, &self.expression, &[])?)?;
        let branch_name = ctx.current_branch()?;

        match expression.is_match(&branch_name) {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Branch name must match pattern: {}", expression),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use assertor::{assert_that, EqualityAssertion};
    use std::collections::HashMap;

    #[test]
    fn test_branch_name_regex() -> anyhow::Result<()> {
        let rule = BranchNameRegex::new("branch_name_regex".to_string(), t!(r"^feat/.*-feature$"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Success { name, .. } = rule.check(&ctx)? else {
            panic!()
        };

        assert_that!(name).is_equal_to("branch_name_regex".to_string());

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_failure() -> anyhow::Result<()> {
        let rule = BranchNameRegex::new("branch_name_regex".to_string(), t!(r"^feat/.*-bugfix$"));
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let RuleResult::Failure { name, message } = rule.check(&ctx)? else {
            panic!()
        };

        assert_that!(name).is_equal_to("branch_name_regex".to_string());
        assert_that!(message)
            .is_equal_to("Branch name must match pattern: ^feat/.*-bugfix$".to_string());

        Ok(())
    }

    #[test]
    fn test_sync() {
        let rule = BranchNameRegex::new("branch_name_regex".to_string(), t!(r"^feat/.*$"));
        assert!(rule.sync());
    }
}

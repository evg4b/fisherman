use crate::context::Context;
use crate::rules::helpers::match_expression;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;

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
        match match_expression(&self.expression, &ctx.current_branch()?)? {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: String::new(),
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Branch name does not match regex: {}", self.name),
            }),
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
    fn test_branch_name_regex() -> anyhow::Result<()> {
        let rule = BranchNameRegex::new(
            "branch_name_regex".to_string(),
            tmpl!(r"^feat/.*-feature$"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("feat/my-feature".to_string()));

        let RuleResult::Success { name, .. } = rule.check(&ctx)? else { panic!() };

        assert_that!(name).is_equal_to("branch_name_regex".to_string());

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_failure() -> anyhow::Result<()> {
        let rule = BranchNameRegex::new(
            "branch_name_regex".to_string(),
            tmpl!(r"^feat/.*-bugfix$"),
        );
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("bugfix/my-feature".to_string()));

        let RuleResult::Failure { name, message } = rule.check(&ctx)? else { panic!() };

        assert_that!(name).is_equal_to("branch_name_regex".to_string());
        assert_that!(message).is_equal_to("Branch name does not match regex: branch_name_regex".to_string());

        Ok(())
    }
}
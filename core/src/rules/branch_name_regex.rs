use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{Rule, RuleResult, ConditionalRule};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use regex::Regex;
use rules_derive::ConditionalRule;
use serde::{Deserialize, Serialize};

static BRANCH_NAME_REGEX_RULE_NAME: &str = "branch-name-regex";

#[derive(Debug, Deserialize, Serialize, ConditionalRule)]
pub struct BranchNameRegexRule {
    pub when: Option<Expression>,
    #[serde(alias = "regex")]
    pub expression: TemplateString,
}

impl std::fmt::Display for BranchNameRegexRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Branch name must match pattern: {}", self.expression)
    }
}

#[typetag::serde(name = "branch-name-regex")]
impl Rule for BranchNameRegexRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: BRANCH_NAME_REGEX_RULE_NAME.to_string(),
            });
        }

        let expression = Regex::new(&compile_tmpl(ctx, &self.expression, &[])?)?;
        let branch_name = ctx.current_branch()?;

        match expression.is_match(&branch_name) {
            true => Ok(RuleResult::Success {
                name: BRANCH_NAME_REGEX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: BRANCH_NAME_REGEX_RULE_NAME.to_string(),
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
    use assert2::assert;
    use std::collections::HashMap;

    #[test]
    fn test_branch_name_regex_success() -> anyhow::Result<()> {
        let rule = BranchNameRegexRule {
            when: None,
            expression: t!(r"^feat/.*-feature$"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Success { name, .. } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "branch-name-regex");

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_failure() -> anyhow::Result<()> {
        let rule = BranchNameRegexRule {
            when: None,
            expression: t!(r"^feat/.*-bugfix$"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx)?;
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "branch-name-regex");
        assert_eq!(message, "Branch name must match pattern: ^feat/.*-bugfix$");

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_variables_error() -> anyhow::Result<()> {
        let rule = BranchNameRegexRule {
            when: None,
            expression: t!(r"^feat/.*$"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/test".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_branch_error() -> anyhow::Result<()> {
        let rule = BranchNameRegexRule {
            when: None,
            expression: t!(r"^feat/.*$"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Err(anyhow::anyhow!("Branch error")));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_branch_name_regex_invalid_regex() -> anyhow::Result<()> {
        let rule = BranchNameRegexRule {
            when: None,
            expression: t!(r"^feat/["),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/test".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_display() {
        let rule = BranchNameRegexRule { when: None, expression: "^feat/".into() };
        assert_eq!(format!("{}", rule), "Branch name must match pattern: `^feat/`");
    }
}

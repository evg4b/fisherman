use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::rule::{Rule, RuleResult, ConditionalRule};
use crate::rules::CompiledRule;
use crate::templates::TemplateString;
use crate::scripting::Expression;
use rules_derive::ConditionalRule as ConditionalRuleDerive;

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct BranchNameSuffixRule {
    pub when: Option<Expression>,
    pub suffix: TemplateString,
}

static BRANCH_NAME_SUFFIX_RULE_NAME: &str = "branch-name-suffix";

#[typetag::serde(name = "branch-name-suffix")]
impl Rule for BranchNameSuffixRule {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<crate::rules::rule::RuleResult> {
        if self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: BRANCH_NAME_SUFFIX_RULE_NAME.to_string(),
            });
        }

        let suffix = compile_tmpl(ctx, &self.suffix, &[])?;
        let branch_name = ctx.current_branch()?;

        match branch_name.ends_with(&suffix) {
            true => Ok(RuleResult::Success {
                name: BRANCH_NAME_SUFFIX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: BRANCH_NAME_SUFFIX_RULE_NAME.to_string(),
                message: format!("Branch name must end with: {}", suffix),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_branch_name_suffix_success() -> anyhow::Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = BranchNameSuffixRule {
            when: None,
            suffix: t!("feature"),
        }
        .check(&ctx)?;

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

        let result = BranchNameSuffixRule {
            when: None,
            suffix: t!("suffix"),
        }
        .check(&ctx)?;

        assert!(matches!(result, RuleResult::Failure { .. }));

        Ok(())
    }

    #[test]
    fn test_branch_name_suffix_variables_error() {
        let rule = BranchNameSuffixRule {
            when: None,
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("my-suffix".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_branch_name_suffix_branch_error() {
        let rule = BranchNameSuffixRule {
            when: None,
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Err(anyhow::anyhow!("Branch error")));
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }
}

use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{Rule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct BranchNameSuffixRule {
    pub suffix: TemplateString,
}

impl std::fmt::Display for BranchNameSuffixRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Branch name must end with: {}", self.suffix)
    }
}

static BRANCH_NAME_SUFFIX_RULE_NAME: &str = "branch-name-suffix";

#[typetag::serde(name = "branch-name-suffix")]
impl Rule for BranchNameSuffixRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
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
    use anyhow::{anyhow, Result};
    use std::collections::HashMap;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = BranchNameSuffixRule {
            suffix: t!("-patch"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"suffix":"-patch"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: BranchNameSuffixRule = serde_json::from_str(r#"{"suffix":"-patch"}"#)?;

        assert_eq!(config.suffix, t!("-patch"));

        Ok(())
    }


    #[test]
    fn test_branch_name_suffix_success() -> Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = BranchNameSuffixRule {
            suffix: t!("feature"),
        }
            .check(&ctx)?;

        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn test_branch_name_suffix_failure() -> Result<()> {
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("bugfix/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = BranchNameSuffixRule {
            suffix: t!("suffix"),
        }
            .check(&ctx)?;

        assert!(matches!(result, RuleResult::Failure { .. }));

        Ok(())
    }

    #[test]
    fn test_branch_name_suffix_variables_error() {
        let rule = BranchNameSuffixRule {
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("my-suffix".to_string()));
        ctx.expect_variables()
            .returning(|| Err(anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_branch_name_suffix_branch_error() {
        let rule = BranchNameSuffixRule {
            suffix: t!("suffix"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Err(anyhow!("Branch error")));
        ctx.expect_variables()
            .returning(|| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = BranchNameSuffixRule {
            suffix: "-patch".into(),
        };
        assert_eq!(format!("{}", rule), "Branch name must end with: `-patch`");
    }
}

use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fmt::Display;

static BRANCH_NAME_PREFIX_RULE_NAME: &str = "branch-name-prefix";

#[derive(Debug, Deserialize, Serialize)]
pub struct BranchNamePrefixRule {
    pub prefix: TemplateString,
}

impl Display for BranchNamePrefixRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Branch should start with: {}", self.prefix)
    }
}

#[typetag::serde(name = "branch-name-prefix")]
impl Rule for BranchNamePrefixRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let prefix = compile_tmpl(ctx, &self.prefix, &[])?;
        let branch_name = ctx.current_branch()?;

        match branch_name.starts_with(&prefix) {
            true => Ok(RuleResult::Success {
                name: BRANCH_NAME_PREFIX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: BRANCH_NAME_PREFIX_RULE_NAME.to_string(),
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
    fn serialize_test() -> Result<()> {
        let config = BranchNamePrefixRule {
            prefix: t!("feat/"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"when":null,"prefix":"feat/"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: BranchNamePrefixRule = serde_json::from_str(r#"{"prefix":"feat/"}"#)?;

        assert!(config.when.is_none());
        assert_eq!(config.prefix, t!("feat/"));

        Ok(())
    }

    #[test]
    fn serialize_test_with_when() -> Result<()> {
        let config = BranchNamePrefixRule {
            when: Some(Expression::new("is_def_var(\"ticket\")")),
            prefix: t!("feat/"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"when":"is_def_var(\"ticket\")","prefix":"feat/"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test_with_when() -> Result<()> {
        let config: BranchNamePrefixRule = serde_json::from_str(
            r#"{"when":"is_def_var(\"ticket\")","prefix":"feat/"}"#,
        )?;

        assert!(config.when.is_some());
        assert_eq!(config.prefix, t!("feat/"));

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_success() -> Result<()> {
        let rule = BranchNamePrefixRule {
            when: None,
            prefix: t!("feat/"),
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
        assert_eq!(name, "branch-name-prefix");

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_failure() -> anyhow::Result<()> {
        let rule = BranchNamePrefixRule {
            when: None,
            prefix: t!("feat/"),
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
        assert_eq!(name, "branch-name-prefix");
        assert_eq!(message, "Branch name must start with: feat/");

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_variables_error() -> anyhow::Result<()> {
        let rule = BranchNamePrefixRule {
            when: None,
            prefix: t!("feat/"),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch()
            .returning(|| Ok("feat/my-feature".to_string()));
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_branch_name_prefix_branch_error() -> anyhow::Result<()> {
        let rule = BranchNamePrefixRule {
            when: None,
            prefix: t!("feat/"),
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
    fn test_display() {
        let rule = BranchNamePrefixRule {
            when: None,
            prefix: "feat/".into(),
        };
        assert_eq!(format!("{}", rule), "Branch should start with: `feat/`");
    }
}

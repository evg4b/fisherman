use crate::context::Context;
use crate::Expression;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fmt::Display;

#[derive(Debug)]
pub enum RuleResult {
    Success {
        name: String,
        output: Option<String>,
    },
    Failure {
        name: String,
        message: String,
    },
    Skipped {
        name: String,
    },
}

#[typetag::serde(tag = "type")]
pub trait Rule: Send + Sync + Display {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult>;
}

#[derive(Serialize, Deserialize)]
pub struct RuleContext {
    pub extract: Option<Vec<String>>,
    pub when: Option<Expression>,
    #[serde(flatten)]
    pub rule: Box<dyn Rule>,
}

impl RuleContext {
    pub fn check_rule(&self, ctx: &mut dyn Context) -> Result<RuleResult> {
        let extended = self.extract.as_ref().map(|e| ctx.extend(e)).transpose()?;
        let correct_ctx: &dyn Context = match &extended {
            Some(boxed) => boxed.as_ref(),
            None => ctx,
        };

        if self.when.is_some() && !self.check_condition(correct_ctx)? {
            return Ok(RuleResult::Skipped {
                name: self.rule.typetag_name().into(),
            });
        }

        self.rule.check(correct_ctx)
    }

    fn check_condition(&self, ctx: &dyn Context) -> Result<bool> {
        self.when
            .as_ref()
            .map(|expr| expr.check(ctx.variables()))
            .unwrap_or(Ok(false))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::{Context, MockContext};
    use crate::rules::BranchNamePrefixRule;
    use crate::scripting::Expression;
    use crate::t;
    use std::collections::HashMap;
    use anyhow::anyhow;

    #[test]
    fn test_deserialize() {
        let json = r#"{"loglogol":null,"type":"branch-name-prefix","prefix":"feat:"}"#;
        let rule: RuleContext = serde_json::from_str(json).unwrap();
        assert_eq!(rule.extract, None);
        assert_eq!(rule.rule.typetag_name(), "branch-name-prefix");
    }

    #[test]
    fn test_deserialize_with_extract() -> Result<()> {
        let json = r#"{
            "extract":["branch:^(?P<Type>feature|bugfix)"],
            "type":"branch-name-prefix",
            "prefix":"feat:"
        }"#;
        let rule: RuleContext = serde_json::from_str(json)?;
        assert_eq!(rule.extract, Some(vec!["branch:^(?P<Type>feature|bugfix)".into()]));

        Ok(())
    }

    #[test]
    fn test_deserialize_with_when() -> Result<()> {
        let json = r#"{
            "extract":["branch:^(?P<Type>feature|bugfix)"],
            "type":"branch-name-prefix",
            "prefix":"feat:",
            "when": "branch.startsWith('feat')"
        }"#;

        let rule: RuleContext = serde_json::from_str(json)?;

        assert_eq!(rule.extract, Some(vec!["branch:^(?P<Type>feature|bugfix)".into()]));
        Ok(())
    }

    #[test]
    fn check_rule_no_extract_no_when_success() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: None,
            when: None,
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("feat/something".to_string()));
        ctx.expect_variables().returning(|| Ok(HashMap::new()));

        let result = rule_ctx.check_rule(&mut ctx)?;
        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn check_rule_no_extract_no_when_failure() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: None,
            when: None,
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_current_branch().returning(|| Ok("bugfix/something".to_string()));
        ctx.expect_variables().returning(|| Ok(HashMap::new()));

        let result = rule_ctx.check_rule(&mut ctx)?;
        assert!(matches!(result, RuleResult::Failure { .. }));

        Ok(())
    }

    #[test]
    fn check_rule_with_extract_extends_context() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: Some(vec!["branch:^(?P<Type>feat|fix)".to_string()]),
            when: None,
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_extend().returning(|_| {
            let mut inner = MockContext::new();
            inner.expect_current_branch().returning(|| Ok("feat/something".to_string()));
            inner.expect_variables().returning(|| Ok(HashMap::new()));
            Ok(Box::new(inner) as Box<dyn Context>)
        });

        let result = rule_ctx.check_rule(&mut ctx)?;
        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn check_rule_with_when_false_returns_skipped() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: None,
            when: Some(Expression::new("1 < 0")),
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(|| Ok(HashMap::new()));

        let result = rule_ctx.check_rule(&mut ctx)?;
        assert!(matches!(result, RuleResult::Skipped { .. }));

        Ok(())
    }

    #[test]
    fn check_rule_with_when_true_runs_rule() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: None,
            when: Some(Expression::new("1 > 0")),
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(|| Ok(HashMap::new()));
        ctx.expect_current_branch().returning(|| Ok("feat/something".to_string()));

        let result = rule_ctx.check_rule(&mut ctx)?;
        assert!(matches!(result, RuleResult::Success { .. }));

        Ok(())
    }

    #[test]
    fn check_rule_condition_error_propagates() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: None,
            when: Some(Expression::new("1 > 0")),
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(|| Err(anyhow!("variables error")));

        let result = rule_ctx.check_rule(&mut ctx);
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn check_rule_extend_error_propagates() -> Result<()> {
        let rule_ctx = RuleContext {
            extract: Some(vec!["branch:something".to_string()]),
            when: None,
            rule: Box::new(BranchNamePrefixRule { prefix: t!("feat/") }),
        };
        let mut ctx = MockContext::new();
        ctx.expect_extend().returning(|_| Err(anyhow!("extend error")));

        let result = rule_ctx.check_rule(&mut ctx);
        assert!(result.is_err());

        Ok(())
    }
}

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
}

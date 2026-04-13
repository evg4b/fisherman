use crate::context::Context;
use crate::Expression;
use anyhow::Result;
use serde::{Deserialize, Serialize};
use std::fmt::Display;

/// Determines how a rule is scheduled during hook execution.
///
/// Sync rules run sequentially in a single thread. Async rules run
/// concurrently in a thread pool (bounded by [`MAX_CONCURRENT_ASYNC_RULES`]).
/// Both groups execute in parallel with each other.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ExecutionMode {
    /// Run sequentially in the caller's thread, one after another.
    Sync,
    /// Run concurrently in a dedicated thread, limited by the global semaphore.
    Async,
}

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

    /// Returns how this rule should be scheduled at execution time.
    ///
    /// The default is [`ExecutionMode::Sync`]. Override to return
    /// [`ExecutionMode::Async`] for rules that perform file I/O or
    /// run external processes and can safely run in parallel.
    fn execution_mode(&self) -> ExecutionMode {
        ExecutionMode::Sync
    }
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
        let extended = self.extract.as_ref()
            .map(|e| ctx.extend(e))
            .transpose()?;

        let actual_ctx: &dyn Context = match extended.as_ref() {
            Some(boxed) => boxed.as_ref(),
            None => ctx,
        };

        if self.when.is_some() && !self.check_condition(actual_ctx)? {
            return Ok(RuleResult::Skipped {
                name: self.rule.typetag_name().into(),
            });
        }

        self.rule.check(actual_ctx)
    }

    fn check_condition(&self, ctx: &dyn Context) -> Result<bool> {
        self.when.as_ref()
            .map(|expr| expr.check(ctx.variables()))
            .unwrap_or(Ok(false))
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::{Context, MockContext};
    use crate::rules::{
        BranchNamePrefixRule, BranchNameRegexRule, BranchNameSuffixRule,
        CommitMessagePrefixRule, CommitMessageRegexRule, CommitMessageSuffixRule,
        CopyFilesRule, DeleteFilesRule, ExecRule, ShellScriptRule,
        SuppressFilesRule, SuppressStringRule, WriteFileRule,
    };
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

    // ── ExecutionMode default and overrides ───────────────────────────────────

    #[test]
    fn sync_rules_default_to_sync_mode() {
        assert_eq!(
            BranchNamePrefixRule { prefix: t!("feat/") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            BranchNameSuffixRule { suffix: t!("-ok") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            BranchNameRegexRule { expression: t!(".*") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            CommitMessagePrefixRule { prefix: t!("fix:") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            CommitMessageSuffixRule { suffix: t!(".") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            CommitMessageRegexRule { expression: t!(".*"), when: None }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            SuppressFilesRule { glob: t!("*.log") }.execution_mode(),
            ExecutionMode::Sync
        );
        assert_eq!(
            SuppressStringRule { regex: t!("TODO"), glob: None }.execution_mode(),
            ExecutionMode::Sync
        );
    }

    #[test]
    fn async_rules_return_async_mode() {
        assert_eq!(
            ExecRule { command: "echo".into(), args: None, env: None }.execution_mode(),
            ExecutionMode::Async
        );
        assert_eq!(
            ShellScriptRule { script: t!("echo hi"), env: None }.execution_mode(),
            ExecutionMode::Async
        );
        assert_eq!(
            WriteFileRule {
                path: t!("out.txt"),
                content: t!("hello"),
                append: None,
            }
            .execution_mode(),
            ExecutionMode::Async
        );
        assert_eq!(
            CopyFilesRule {
                glob: t!("*.txt"),
                src: None,
                destination: t!("dst/"),
            }
            .execution_mode(),
            ExecutionMode::Async
        );
        assert_eq!(
            DeleteFilesRule {
                glob: t!("*.tmp"),
                fail_if_not_found: false,
            }
            .execution_mode(),
            ExecutionMode::Async
        );
    }
}

use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use regex::Regex;

static MESSAGE_REGEX_RULE_NAME: &str = "message-regex";

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct CommitMessageRegexRule {
    pub when: Option<Expression>,
    #[serde(alias = "regex")]
    pub expression: TemplateString,
}

impl std::fmt::Display for CommitMessageRegexRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Commit message must match pattern: {}", self.expression)
    }
}

#[typetag::serde(name = "message-regex")]
impl Rule for CommitMessageRegexRule {
    fn check(&self, ctx: &dyn Context) -> anyhow::Result<RuleResult> {
        let expression = Regex::new(&compile_tmpl(ctx, &self.expression, &[])?)?;
        let commit_msg = ctx.commit_msg()?;

        match expression.is_match(&commit_msg) {
            true => Ok(RuleResult::Success {
                name: MESSAGE_REGEX_RULE_NAME.to_string(),
                output: None,
            }),
            false => Ok(RuleResult::Failure {
                name: MESSAGE_REGEX_RULE_NAME.to_string(),
                message: format!("Commit message must match pattern: {}", expression),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::context::MockContext;
    use crate::rules::CommitMessageRegexRule;
    use crate::rules::{Rule, RuleResult};
    use crate::scripting::Expression;
    use crate::t;
    use anyhow::Result;
    use std::collections::HashMap;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = CommitMessageRegexRule {
            when: None,
            expression: t!(r"^feat:"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"when":null,"expression":"^feat:"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: CommitMessageRegexRule = serde_json::from_str(r#"{"expression":"^feat:"}"#)?;

        assert!(config.when.is_none());
        assert_eq!(config.expression, t!(r"^feat:"));

        Ok(())
    }

    #[test]
    fn serialize_test_with_when() -> Result<()> {
        let config = CommitMessageRegexRule {
            when: Some(Expression::new("is_def_var(\"Ticket\")")),
            expression: t!(r"^feat:"),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(serialized, r#"{"when":"is_def_var(\"Ticket\")","expression":"^feat:"}"#);

        Ok(())
    }

    #[test]
    fn deserialize_test_with_when() -> Result<()> {
        let config: CommitMessageRegexRule = serde_json::from_str(
            r#"{"when":"is_def_var(\"Ticket\")","expression":"^feat:"}"#,
        )?;

        assert!(config.when.is_some());
        assert_eq!(config.expression, t!(r"^feat:"));

        Ok(())
    }

    #[test]
    fn deserialize_test_with_regex_alias() -> Result<()> {
        let config: CommitMessageRegexRule = serde_json::from_str(r#"{"regex":"^feat:"}"#)?;

        assert!(config.when.is_none());
        assert_eq!(config.expression, t!(r"^feat:"));

        Ok(())
    }

    #[test]
    fn deserialize_test_with_regex_alias_and_when() -> Result<()> {
        let config: CommitMessageRegexRule = serde_json::from_str(
            r#"{"when":"is_def_var(\"Ticket\")","regex":"^feat:"}"#,
        )?;

        assert!(config.when.is_some());
        assert_eq!(config.expression, t!(r"^feat:"));

        Ok(())
    }

    #[test]
    fn test_commit_message_regex() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^Test"),
        };
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Test commit message".to_string()));
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));
        let result = rule.check(&context).unwrap();
        match result {
            RuleResult::Success { name, output } => {
                assert_eq!(name, "message-regex");
                assert_eq!(output, None);
            }
            RuleResult::Failure { name, message } => {
                panic!("Expected success, got failure: {} - {}", name, message);
            }
            RuleResult::Skipped { name } => {
                panic!("Expected success, got skipped: {}", name);
            }
        }
    }

    #[test]
    fn test_commit_message_regex_failure() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^Test"),
        };
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Invalid commit message".to_string()));
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));
        let result = rule.check(&context).unwrap();
        match result {
            RuleResult::Success { name, output } => {
                panic!("Expected failure, got success: {} - {:?}", name, output);
            }
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "message-regex");
                assert_eq!(message, "Commit message must match pattern: ^Test");
            }
            RuleResult::Skipped { name } => {
                panic!("Expected failure, got skipped: {}", name);
            }
        }
    }

    #[test]
    fn test_commit_message_regex_error() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^Test"),
        };
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Err(anyhow::anyhow!("Error")));
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));
        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_commit_message_regex_variables_error() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^Test"),
        };
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Test message".to_string()));
        context
            .expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));
        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_commit_message_regex_invalid_regex() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: t!("^Test["),
        };
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Test message".to_string()));
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));
        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = CommitMessageRegexRule {
            when: None,
            expression: "^feat:".into(),
        };
        assert_eq!(
            format!("{}", rule),
            "Commit message must match pattern: `^feat:`"
        );
    }
}

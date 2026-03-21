use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::rule::{Rule, RuleResult};
use crate::rules::{CompiledRule, RuleResultOld};
use crate::templates::TemplateString;
use regex::Regex;

static MESSAGE_REGEX_RULE_NAME: &str = "message-regex";

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct CommitMessageRegexRule {
    #[serde(alias = "regex")]
    pub expression: TemplateString,
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
    use crate::rules::commit_message_regex::CommitMessageRegexRule;
    use crate::rules::rule::{Rule, RuleResult};
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_commit_message_regex() {
        let rule = CommitMessageRegexRule {
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
}

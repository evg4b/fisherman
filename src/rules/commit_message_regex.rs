use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use regex::Regex;

#[derive(Debug)]
pub struct CommitMessageRegex {
    name: String,
    expression: TemplateString,
}

impl CommitMessageRegex {
    pub fn new(name: String, expression: TemplateString) -> Self {
        Self { name, expression }
    }
}

impl CompiledRule for CommitMessageRegex {
    fn check(&self, context: &dyn Context) -> anyhow::Result<RuleResult> {
        let commit_message = context.commit_msg()?;
        let filled_expression = self.expression.to_string()?;
        let expression = Regex::new(&filled_expression)?;

        if expression.is_match(&commit_message) {
            Ok(RuleResult::Success {
                name: self.name.clone(),
                output: String::new(),
            })
        } else {
            Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("Commit message does not match regex: {}", self.name),
            })
        }
    }
}

#[cfg(test)]
mod test {
    use crate::context::MockContext;
    use crate::rules::commit_message_regex::CommitMessageRegex;
    use crate::rules::CompiledRule;
    use crate::rules::RuleResult;

    use crate::tmpl;

    #[test]
    fn test_commit_message_regex() {
        let rule = CommitMessageRegex::new(
            "Test".to_string(),
            tmpl!("^Test"),
        );
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Test commit message".to_string()));
        let result = rule.check(&context).unwrap();
        match result {
            RuleResult::Success { name, output } => {
                assert_eq!(name, "Test");
                assert_eq!(output, "");
            }
            RuleResult::Failure { name, message } => {
                panic!("Expected success, got failure: {} - {}", name, message);
            }
        }
    }

    #[test]
    fn test_commit_message_regex_failure() {
        let rule = CommitMessageRegex::new(
            "Test".to_string(),
            tmpl!("^Test"),
        );
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Ok("Invalid commit message".to_string()));
        let result = rule.check(&context).unwrap();
        match result {
            RuleResult::Success { name, output } => {
                panic!("Expected failure, got success: {} - {}", name, output);
            }
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "Test");
                assert_eq!(message, "Commit message does not match regex: Test");
            }
        }
    }

    #[test]
    fn test_commit_message_regex_error() {
        let rule = CommitMessageRegex::new(
            "Test".to_string(),
            tmpl!("^Test"),
        );
        let mut context = MockContext::new();
        context
            .expect_commit_msg()
            .returning(|| Err(anyhow::anyhow!("Error")));
        let result = rule.check(&context);
        assert!(result.is_err());
    }
}

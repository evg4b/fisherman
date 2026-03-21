use crate::context::Context;
use crate::extract_vars;
use crate::rules::rule::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use rules_derive::ConditionalRule as ConditionalRuleDerive;
use serde::{Deserialize, Serialize};

static SUPPRESS_FILES_RULE_NAME: &str = "suppress-files";

#[derive(Debug, Serialize, Deserialize, ConditionalRuleDerive)]
pub struct SuppressFilesRule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    pub glob: TemplateString,
}

impl std::fmt::Display for SuppressFilesRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Suppress files matching: {}", self.glob)
    }
}

#[typetag::serde(name = "suppress-files")]
impl Rule for SuppressFilesRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: SUPPRESS_FILES_RULE_NAME.to_string(),
            });
        }

        let variables = extract_vars!(self, ctx)?;
        let glob_pattern = self.glob.compile(&variables)?;
        let pattern = Pattern::new(glob_pattern.as_str())?;
        let staged_files = ctx.staged_files()?;

        let mut matched_files = Vec::new();
        for file in staged_files {
            if pattern.matches_path(&file) {
                matched_files.push(file.display().to_string());
            }
        }

        if !matched_files.is_empty() {
            return Ok(RuleResult::Failure {
                name: SUPPRESS_FILES_RULE_NAME.to_string(),
                message: format!(
                    "The following files are suppressed from being committed: {}",
                    matched_files.join(", ")
                ),
            });
        }

        Ok(RuleResult::Success {
            name: SUPPRESS_FILES_RULE_NAME.to_string(),
            output: None,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;
    use std::path::PathBuf;

    #[test]
    fn test_suppress_files_success() -> Result<()> {
        let rule = SuppressFilesRule {
            when: None,
            glob: tmpl!("*.txt"),
            extract: None,
        };
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context.expect_staged_files().returning(|| {
            Ok(vec![
                PathBuf::from("README.md"),
                PathBuf::from("src/main.rs"),
            ])
        });

        let result = rule.check(&context)?;
        match result {
            RuleResult::Success { .. } => {}
            _ => panic!("Expected Success"),
        }
        Ok(())
    }

    #[test]
    fn test_suppress_files_failure() -> Result<()> {
        let rule = SuppressFilesRule {
            when: None,
            extract: None,
            glob: tmpl!("*.txt"),
        };
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context.expect_staged_files().returning(|| {
            Ok(vec![
                PathBuf::from("test.txt"),
                PathBuf::from("src/main.rs"),
            ])
        });

        let result = rule.check(&context)?;
        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "suppress-files");
                assert!(
                    message.contains(
                        "The following files are suppressed from being committed: test.txt"
                    )
                );
            }
            _ => panic!("Expected failure"),
        }
        Ok(())
    }

    #[test]
    fn test_display() {
        let rule = SuppressFilesRule { when: None, extract: None, glob: "*.secret".into() };
        assert_eq!(format!("{}", rule), "Suppress files matching: `*.secret`");
    }
}

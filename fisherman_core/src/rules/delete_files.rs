use crate::context::Context;
use crate::rules::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::{bail, Result};
use glob::{glob, GlobResult};
use rules_derive::ConditionalRule as ConditionalRuleDerive;
use std::fs;

static DELETE_FILES_RULE_NAME: &str = "delete-files";

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct DeleteFilesRule {
    pub when: Option<Expression>,
    pub glob: TemplateString,
    pub fail_if_not_found: bool,
}

impl std::fmt::Display for DeleteFilesRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Delete files matching '{}'", self.glob)
    }
}

#[typetag::serde(name = "delete-files")]
impl Rule for DeleteFilesRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: DELETE_FILES_RULE_NAME.to_string(),
            });
        }

        let variables = ctx.variables(&[])?;
        let glob_pattern = self.glob.compile(&variables)?;
        let paths = glob(glob_pattern.as_str())?.collect::<Vec<GlobResult>>();

        if paths.is_empty() && self.fail_if_not_found {
            return Ok(RuleResult::Failure {
                name: DELETE_FILES_RULE_NAME.to_string(),
                message: format!("No files matched the glob pattern: {}", glob_pattern),
            });
        }

        for path in paths {
            match path {
                Ok(path) => fs::remove_file(path.as_path())?,
                Err(err) => {
                    bail!("Error deleting file: {}", err);
                }
            }
        }

        Ok(RuleResult::Success {
            name: DELETE_FILES_RULE_NAME.to_string(),
            output: None,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;
    use std::fs::File;
    use tempfile::tempdir;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = DeleteFilesRule {
            when: None,
            glob: "*.log".into(),
            fail_if_not_found: false,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"when":null,"glob":"*.log","fail_if_not_found":false}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: DeleteFilesRule =
            serde_json::from_str(r#"{"glob":"*.log","fail_if_not_found":true}"#)?;

        assert!(config.when.is_none());
        assert_eq!(config.glob, "*.log".into());
        assert_eq!(config.fail_if_not_found, true);

        Ok(())
    }

    #[test]
    fn test_delete_files_success() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path = temp_dir.path().join("test_file.txt");

        File::create(&file_path)?;

        let rule = DeleteFilesRule {
            when: None,
            glob: tmpl!(file_path.display()),
            fail_if_not_found: true,
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        let result = rule.check(&context)?;

        match result {
            RuleResult::Success { .. } => {}
            _ => panic!("Expected Success"),
        }
        assert!(!file_path.exists());

        Ok(())
    }

    #[test]
    fn test_delete_files_no_matches_with_failure() -> Result<()> {
        let rule = DeleteFilesRule {
            when: None,
            glob: tmpl!("path/that/does/not/exist/*.txt"),
            fail_if_not_found: true,
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        let result = rule.check(&context)?;

        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "delete-files");
                assert!(message.contains("No files matched the glob pattern"));
            }
            _ => panic!("Expected RuleResult::Failure"),
        }

        Ok(())
    }

    #[test]
    fn test_delete_files_no_matches_without_failure() -> Result<()> {
        let rule = DeleteFilesRule {
            when: None,
            glob: tmpl!("path/that/does/not/exist/*.txt"),
            fail_if_not_found: false,
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context)?;

        match result {
            RuleResult::Success { .. } => {}
            _ => panic!("Expected Success"),
        }

        Ok(())
    }

    #[test]
    fn test_delete_multiple_files() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path1 = temp_dir.path().join("test_file1.txt");
        let file_path2 = temp_dir.path().join("test_file2.txt");

        File::create(&file_path1)?;
        File::create(&file_path2)?;

        let glob_pattern = format!("{}/*.txt", temp_dir.path().display());
        let rule = DeleteFilesRule {
            when: None,
            glob: tmpl!(glob_pattern),
            fail_if_not_found: true,
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context)?;

        match result {
            RuleResult::Success { .. } => {}
            _ => panic!("Expected Success"),
        }
        assert!(!file_path1.exists());
        assert!(!file_path2.exists());

        Ok(())
    }

    #[test]
    fn test_delete_files_glob_error() {
        let rule = DeleteFilesRule {
            when: None,
            glob: tmpl!("[invalid-glob"),
            fail_if_not_found: true,
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = DeleteFilesRule {
            when: None,
            glob: "*.log".into(),
            fail_if_not_found: false,
        };
        assert_eq!(format!("{}", rule), "Delete files matching '`*.log`'");
    }
}

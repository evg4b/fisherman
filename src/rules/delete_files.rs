use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::{bail, Result};
use glob::{glob, GlobResult};
use std::fs;

pub struct DeleteFiles {
    name: String,
    glob: TemplateString,
    fail_if_not_found: bool,
}

impl DeleteFiles {
    pub fn new(name: String, glob: TemplateString, fail_if_not_found: bool) -> DeleteFiles {
        DeleteFiles {
            name,
            glob,
            fail_if_not_found,
        }
    }
}

impl CompiledRule for DeleteFiles {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables(&vec![])?;
        let glob_pattern = self.glob.to_string(&variables)?;
        let paths = glob(glob_pattern.as_str())?
            .collect::<Vec<GlobResult>>();

        if paths.is_empty() && self.fail_if_not_found {
            return Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "No files matched the glob pattern: {}",
                    glob_pattern
                ),
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
            name: self.name.clone(),
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
    fn test_delete_files_success() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path = temp_dir.path().join("test_file.txt");

        File::create(&file_path)?;

        let rule = DeleteFiles::new(
            "test_rule".to_string(),
            tmpl!(file_path.display()),
            true,
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(std::collections::HashMap::new()));
        let result = rule.check(&context)?;

        assert!(matches!(result, RuleResult::Success { .. }));
        assert!(!file_path.exists());

        Ok(())
    }

    #[test]
    fn test_delete_files_no_matches_with_failure() -> Result<()> {
        let rule = DeleteFiles::new(
            String::from("test_rule"),
            tmpl!("path/that/does/not/exist/*.txt"),
            true,
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(std::collections::HashMap::new()));
        let result = rule.check(&context)?;

        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "test_rule");
                assert!(message.contains("No files matched the glob pattern"));
            },
            _ => panic!("Expected RuleResult::Failure"),
        }

        Ok(())
    }

    #[test]
    fn test_delete_files_no_matches_without_failure() -> Result<()> {
        let rule = DeleteFiles::new(
            String::from("test_rule"),
            tmpl!("path/that/does/not/exist/*.txt"),
            false,
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context)?;

        assert!(matches!(result, RuleResult::Success { .. }));

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
        let rule = DeleteFiles::new(
            String::from("test_rule"),
            tmpl!(glob_pattern),
            true,
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context)?;

        assert!(matches!(result, RuleResult::Success { .. }));
        assert!(!file_path1.exists());
        assert!(!file_path2.exists());

        Ok(())
    }

    #[test]
    fn test_delete_files_glob_error() {
        let rule = DeleteFiles::new(
            String::from("test_rule"),
            tmpl!("[invalid-glob"),
            true,
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(std::collections::HashMap::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }
}
use crate::context::Context;
use crate::rules::{Rule, RuleResult};
use crate::templates::TemplateString;
use anyhow::{bail, Result};
use glob::glob;
use std::fs;
use std::fs::create_dir_all;
use std::path::Path;

static COPY_FILES_RULE_NAME: &str = "copy-files";

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct CopyFilesRule {
    pub glob: TemplateString,
    pub src: Option<TemplateString>,
    pub destination: TemplateString,
}

impl std::fmt::Display for CopyFilesRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self.src {
            None => write!(
                f,
                "Copy files matching '{}' to '{}'",
                self.glob, self.destination
            ),
            Some(ref src) => write!(
                f,
                "Copy files matching '{}' from '{}' to '{}'",
                self.glob, src, self.destination
            ),
        }
    }
}

fn ensure_parent_exists(path: &Path) -> Result<()> {
    if let Some(parent) = path.parent()
        && !parent.exists()
    {
        create_dir_all(parent)?;
    }
    Ok(())
}

#[typetag::serde(name = "copy-files")]
impl Rule for CopyFilesRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables()?;
        let compiled_glob = self.glob.compile(&variables)?;
        let compiled_src = self
            .src
            .as_ref()
            .map(|s| s.compile(&variables))
            .transpose()?;
        let compiled_pattern = match compiled_src.clone() {
            Some(src) => Path::join(src.as_ref(), compiled_glob),
            None => compiled_glob.parse()?,
        };
        let compiled_destination = self.destination.compile(&variables)?;

        let mut copied_files = 0;

        for entry in glob(compiled_pattern.to_str().unwrap())? {
            match entry {
                Ok(path) => {
                    let new_name = match compiled_src.as_ref() {
                        Some(value) => path.strip_prefix(value)?.display().to_string(),
                        None => path.display().to_string(),
                    };
                    let destination_path = Path::join(compiled_destination.as_ref(), new_name);

                    ensure_parent_exists(&destination_path)?;
                    fs::copy(&path, &destination_path)?;

                    copied_files += 1;
                }
                Err(e) => {
                    bail!("Error reading glob entry: {}", e);
                }
            }
        }

        Ok(RuleResult::Success {
            name: COPY_FILES_RULE_NAME.to_string(),
            output: Some(format!("Copied {} files", copied_files)),
        })
    }
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;
    use super::*;
    use crate::context::MockContext;
    use crate::tmpl;
    use anyhow::Result;
    use assertor::{assert_that, EqualityAssertion};
    use std::env;
    use std::fs::File;
    use std::io::Write;
    use tempfile::tempdir;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = CopyFilesRule {
            glob: "*.txt".into(),
            src: Some("src/".into()),
            destination: "dist/".into(),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"glob":"*.txt","src":"src/","destination":"dist/"}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: CopyFilesRule = serde_json::from_str(
            r#"{"glob":"*.txt","src":"src/","destination":"dist/"}"#,
        )?;

        assert_eq!(config.glob, "*.txt".into());
        assert_eq!(config.src, Some("src/".into()));
        assert_eq!(config.destination, "dist/".into());

        Ok(())
    }

    #[test]
    fn serialize_test_without_src() -> Result<()> {
        let config = CopyFilesRule {
            glob: "*.txt".into(),
            src: None,
            destination: "dist/".into(),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"glob":"*.txt","src":null,"destination":"dist/"}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test_without_src() -> Result<()> {
        let config: CopyFilesRule = serde_json::from_str(
            r#"{"glob":"*.txt","destination":"dist/"}"#,
        )?;

        assert_eq!(config.glob, "*.txt".into());
        assert!(config.src.is_none());
        assert_eq!(config.destination, "dist/".into());

        Ok(())
    }


    #[test]
    fn test_copy_files_with_src() -> Result<()> {
        // Create source directory structure
        let temp_src = tempdir()?;
        let src_path = temp_src.path().to_str().unwrap().to_string();

        // Create destination directory
        let temp_dest = tempdir()?;
        let dest_path = temp_dest.path().to_str().unwrap().to_string();

        // Create test file in source directory
        let test_file_path = temp_src.path().join("test.txt");
        let mut file = File::create(&test_file_path)?;
        writeln!(file, "test content")?;

        // Create rule with explicit source
        let rule = CopyFilesRule {
            glob: tmpl!("*.txt".to_string()),
            src: Some(tmpl!(src_path)),
            destination: tmpl!(dest_path),
        };

        // Run the rule
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));
        let result = rule.check(&context)?;
        match result {
            RuleResult::Success { name, output } => {
                assert_that!(name).is_equal_to("copy-files".to_string());
                assert_that!(output).is_equal_to(Some("Copied 1 files".to_string()));
            }
            _ => panic!("Expected Success, but got {:?}", result),
        }

        // Check that file was copied
        let copied_file = std::fs::read_to_string(temp_dest.path().join("test.txt"))?;
        assert_that!(copied_file).is_equal_to("test content\n".to_string());

        Ok(())
    }

    #[test]
    fn test_copy_files_without_src() -> Result<()> {
        // Create temp directory and save current directory
        let current_dir = env::current_dir()?;
        let temp_dir = tempdir()?;
        let temp_dest = tempdir()?;

        // Change to temp directory to create files
        env::set_current_dir(&temp_dir)?;

        // Create test file in current directory
        let test_file_path = Path::new("test-no-src.txt");
        let mut file = File::create(test_file_path)?;
        writeln!(file, "content without src")?;

        // Create rule without source (should use current directory)
        let rule = CopyFilesRule {
            glob: tmpl!("test-no-src.txt".to_string()),
            src: None,
            destination: tmpl!(temp_dest.path().to_str().unwrap().to_string()),
        };

        // Run the rule
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));
        let result = rule.check(&context)?;
        match result {
            RuleResult::Success { name, output } => {
                assert_that!(name).is_equal_to("copy-files".to_string());
                assert_that!(output).is_equal_to(Some("Copied 1 files".to_string()));
            }
            _ => panic!("Expected Success, but got {:?}", result),
        }

        // Check that file was copied
        let copied_file = std::fs::read_to_string(temp_dest.path().join("test-no-src.txt"))?;
        assert_that!(copied_file).is_equal_to("content without src\n".to_string());

        // Return to original directory
        env::set_current_dir(current_dir)?;

        Ok(())
    }

    #[test]
    fn test_ensure_parent_exists() -> Result<()> {
        let temp_dir = tempdir()?;
        let nested_path = temp_dir.path().join("nested").join("dir").join("file.txt");

        ensure_parent_exists(&nested_path)?;

        let parent = nested_path.parent().unwrap();
        assert_that!(parent.exists()).is_equal_to(true);

        Ok(())
    }

    #[test]
    fn test_copy_files_no_matches() -> Result<()> {
        let temp_src = tempdir()?;
        let temp_dest = tempdir()?;

        let rule = CopyFilesRule {
            glob: tmpl!("*.nonexistent".to_string()),
            src: Some(tmpl!(temp_src.path().to_str().unwrap().to_string())),
            destination: tmpl!(temp_dest.path().to_str().unwrap().to_string()),
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));
        let result = rule.check(&context)?;
        match result {
            RuleResult::Success { name, output } => {
                assert_that!(name).is_equal_to("copy-files".to_string());
                assert_that!(output).is_equal_to(Some("Copied 0 files".to_string()));
            }
            _ => panic!("Expected Success, but got {:?}", result),
        }

        Ok(())
    }

    #[test]
    fn test_copy_files_variables_error() {
        let rule = CopyFilesRule {
            glob: tmpl!("*.txt".to_string()),
            src: None,
            destination: tmpl!("/tmp/dest".to_string()),
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_copy_files_glob_template_error() {
        let rule = CopyFilesRule {
            glob: tmpl!("{{missing_var}}/*.txt".to_string()),
            src: None,
            destination: tmpl!("/tmp/dest".to_string()),
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_copy_files_destination_template_error() {
        let rule = CopyFilesRule {
            glob: tmpl!("*.txt".to_string()),
            src: None,
            destination: tmpl!("{{missing_dest}}".to_string()),
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_copy_files_src_template_error() {
        let rule = CopyFilesRule {
            glob: tmpl!("*.txt".to_string()),
            src: Some(tmpl!("{{missing_src}}".to_string())),
            destination: tmpl!("/tmp/dest".to_string()),
        };

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|| Ok(HashMap::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_display_without_src() {
        let rule = CopyFilesRule {
            glob: "*.txt".into(),
            src: None,
            destination: "dist/".into(),
        };
        assert_eq!(
            format!("{}", rule),
            "Copy files matching '`*.txt`' to '`dist/`'"
        );
    }

    #[test]
    fn test_display_with_src() {
        let rule = CopyFilesRule {
            glob: "*.txt".into(),
            src: Some("src/".into()),
            destination: "dist/".into(),
        };
        assert_eq!(
            format!("{}", rule),
            "Copy files matching '`*.txt`' from '`src/`' to '`dist/`'"
        );
    }
}

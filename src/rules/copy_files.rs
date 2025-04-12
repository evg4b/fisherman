use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::{bail, Result};
use glob::glob;
use std::fs;
use std::fs::create_dir_all;
use std::path::Path;

pub struct CopyFiles {
    name: String,
    glob: TemplateString,
    src: Option<TemplateString>,
    destination: TemplateString,
}

impl CopyFiles {
    pub fn new(
        name: String,
        glob: TemplateString,
        src: Option<TemplateString>,
        destination: TemplateString,
    ) -> CopyFiles {
        CopyFiles {
            name,
            glob,
            src,
            destination,
        }
    }

    fn ensure_parent_exists(path: &Path) -> Result<()> {
        if let Some(parent) = path.parent() {
            if !parent.exists() {
                create_dir_all(parent)?;
            }
        }
        Ok(())
    }
}

impl CompiledRule for CopyFiles {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables(&vec![])?;
        let compiled_glob = self.glob.to_string(&variables)?;
        let compiled_src = self.src.as_ref().map(|s| s.to_string(&variables)).transpose()?;
        let compiled_pattern = match compiled_src.clone() {
            Some(src) => Path::join(src.as_ref(), compiled_glob),
            None => compiled_glob.parse()?,
        };
        let compiled_destination = self.destination.to_string(&variables)?;

        let mut copied_files = 0;

        for entry in glob(compiled_pattern.to_str().unwrap())? {
            match entry {
                Ok(path) => {
                    let new_name = match compiled_src.as_ref() {
                        Some(value) => path.strip_prefix(value)?.display().to_string(),
                        None => path.display().to_string(),
                    };
                    let destination_path = Path::join(compiled_destination.as_ref(), new_name);

                    Self::ensure_parent_exists(&destination_path)?;
                    fs::copy(&path, &destination_path)?;

                    copied_files += 1;
                }
                Err(e) => {
                    bail!("Error reading glob entry: {}", e);
                }
            }
        }

        Ok(RuleResult::Success {
            name: self.name.clone(),
            output: Some(format!("Copied {} files", copied_files)),
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use assertor::{assert_that, EqualityAssertion};
    use std::env;
    use std::fs::File;
    use std::io::Write;
    use tempfile::tempdir;
    use crate::tmpl;

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
        let rule = CopyFiles::new(
            "copy-test".to_string(),
            tmpl!("*.txt".to_string()),
            Some(tmpl!(src_path)),
            tmpl!(dest_path),
        );

        // Run the rule
        let result = rule.check(&MockContext::new())?;
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result);
        };

        assert_that!(name).is_equal_to("copy-test".to_string());
        assert_that!(output).is_equal_to("Copied 1 files".to_string());

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
        let rule = CopyFiles::new(
            "copy-test-no-src".to_string(),
            tmpl!("test-no-src.txt".to_string()),
            None,
            tmpl!(temp_dest.path().to_str().unwrap().to_string()),
        );

        // Run the rule
        let result = rule.check(&MockContext::new())?;
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result);
        };

        assert_that!(name).is_equal_to("copy-test-no-src".to_string());
        assert_that!(output).is_equal_to("Copied 1 files".to_string());

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

        CopyFiles::ensure_parent_exists(&nested_path)?;

        let parent = nested_path.parent().unwrap();
        assert_that!(parent.exists()).is_equal_to(true);

        Ok(())
    }

    #[test]
    fn test_copy_files_no_matches() -> Result<()> {
        let temp_src = tempdir()?;
        let temp_dest = tempdir()?;

        let rule = CopyFiles::new(
            "no-matches".to_string(),
            tmpl!("*.nonexistent".to_string()),
            Some(tmpl!(temp_src.path().to_str().unwrap().to_string())),
            tmpl!(temp_dest.path().to_str().unwrap().to_string()),
        );

        let result = rule.check(&MockContext::new())?;
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result);
        };

        assert_that!(name).is_equal_to("no-matches".to_string());
        assert_that!(output).is_equal_to("Copied 0 files".to_string());

        Ok(())
    }
}
use crate::context::Context;
use crate::rules::helpers::compile_tmpl;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use regex::Regex;
use std::fs;

pub struct SuppressString {
    name: String,
    regex: TemplateString,
    glob: Option<TemplateString>,
}

impl SuppressString {
    pub fn new(name: String, regex: TemplateString, glob: Option<TemplateString>) -> Self {
        Self { name, regex, glob }
    }
}

impl CompiledRule for SuppressString {
    fn is_sequential(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables(&[])?;
        let regex_str = compile_tmpl(ctx, &self.regex, &[])?;
        let regex = Regex::new(&regex_str)?;

        let pattern = match &self.glob {
            Some(g) => Some(Pattern::new(&g.to_string(&variables)?)?),
            None => None,
        };

        let staged_files = ctx.staged_files()?;
        let repo_path = ctx.repo_path();

        let mut matched_files = Vec::new();

        for file in staged_files {
            if let Some(p) = &pattern && !p.matches_path(&file) {
                continue;
            }

            let full_path = repo_path.join(&file);
            if !full_path.exists() {
                continue; // Might be a deleted file
            }

            let content = match fs::read_to_string(&full_path) {
                Ok(c) => c,
                Err(_) => continue, // Skip binary files or unreadable files
            };

            if regex.is_match(&content) {
                matched_files.push(file.display().to_string());
            }
        }

        if !matched_files.is_empty() {
            return Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "The following files contain suppressed string: {}",
                    matched_files.join(", ")
                ),
            });
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
    use std::io::Write;
    use tempfile::tempdir;

    #[test]
    fn test_suppress_string_success() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path = temp_dir.path().join("test.txt");
        let mut file = File::create(&file_path)?;
        writeln!(file, "clean content")?;

        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), None);
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_repo_path()
            .return_const(temp_dir.path().to_path_buf());

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }

    #[test]
    fn test_suppress_string_failure() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path = temp_dir.path().join("test.txt");
        let mut file = File::create(&file_path)?;
        writeln!(file, "this has a TODO item")?;

        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), None);
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_repo_path()
            .return_const(temp_dir.path().to_path_buf());

        let result = rule.check(&context)?;
        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "test");
                assert!(message.contains("The following files contain suppressed string: test.txt"));
            }
            _ => panic!("Expected failure"),
        }
        Ok(())
    }

    #[test]
    fn test_suppress_string_with_glob() -> Result<()> {
        let temp_dir = tempdir()?;
        let file_path = temp_dir.path().join("test.txt");
        let mut file = File::create(&file_path)?;
        writeln!(file, "TODO in txt")?;

        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), Some(tmpl!("*.rs")));
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_repo_path()
            .return_const(temp_dir.path().to_path_buf());

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }
}

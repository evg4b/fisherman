use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;

pub struct SuppressFiles {
    name: String,
    glob: TemplateString,
}

impl SuppressFiles {
    pub fn new(name: String, glob: TemplateString) -> Self {
        Self { name, glob }
    }
}

impl CompiledRule for SuppressFiles {
    fn is_sequential(&self) -> bool {
        true
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables(&[])?;
        let glob_pattern = self.glob.to_string(&variables)?;
        let pattern = Pattern::new(&glob_pattern)?;
        let staged_files = ctx.staged_files()?;

        let mut matched_files = Vec::new();
        for file in staged_files {
            if pattern.matches_path(&file) {
                matched_files.push(file.display().to_string());
            }
        }

        if !matched_files.is_empty() {
            return Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!(
                    "The following files are suppressed from being committed: {}",
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
    use std::path::PathBuf;

    #[test]
    fn test_suppress_files_success() -> Result<()> {
        let rule = SuppressFiles::new("test".to_string(), tmpl!("*.txt"));
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![PathBuf::from("README.md"), PathBuf::from("src/main.rs")]));

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }

    #[test]
    fn test_suppress_files_failure() -> Result<()> {
        let rule = SuppressFiles::new("test".to_string(), tmpl!("*.txt"));
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![PathBuf::from("test.txt"), PathBuf::from("src/main.rs")]));

        let result = rule.check(&context)?;
        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "test");
                assert!(message.contains("The following files are suppressed from being committed: test.txt"));
            }
            _ => panic!("Expected failure"),
        }
        Ok(())
    }
}

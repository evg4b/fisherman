use crate::context::Context;
use crate::rules::{CompiledRule, RuleResultOld};
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use serde::{Deserialize, Serialize};
use crate::rules::delete_files::DeleteFilesRule;
use crate::rules::rule::Rule;

#[derive(Debug, Serialize, Deserialize)]
pub struct SuppressFilesRule{
    glob: TemplateString,
}

#[typetag::serde(name = "suppress-files")]
impl Rule for SuppressFilesRule {
    fn check(&self, ctx: &dyn Context) -> Result<crate::rules::rule::RuleResult> {
        todo!()
    }
}


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

    fn check(&self, ctx: &dyn Context) -> Result<RuleResultOld> {
        let variables = ctx.variables(&[])?;
        let glob_pattern = self.glob.compile(&variables)?;
        let pattern = Pattern::new(&glob_pattern)?;
        let staged_files = ctx.staged_files()?;

        let mut matched_files = Vec::new();
        for file in staged_files {
            if pattern.matches_path(&file) {
                matched_files.push(file.display().to_string());
            }
        }

        if !matched_files.is_empty() {
            return Ok(RuleResultOld::Failure {
                name: self.name.clone(),
                message: format!(
                    "The following files are suppressed from being committed: {}",
                    matched_files.join(", ")
                ),
            });
        }

        Ok(RuleResultOld::Success {
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
        assert!(matches!(result, RuleResultOld::Success { .. }));
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
            RuleResultOld::Failure { name, message } => {
                assert_eq!(name, "test");
                assert!(message.contains("The following files are suppressed from being committed: test.txt"));
            }
            _ => panic!("Expected failure"),
        }
        Ok(())
    }
}

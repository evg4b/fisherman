use crate::context::{Context, DiffLine};
use crate::rules::helpers::compile_tmpl;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use regex::Regex;

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
            Some(g) => Some(Pattern::new(&g.compile(&variables)?)?),
            None => None,
        };

        let staged_files = ctx.staged_files()?;

        let mut matched_files = Vec::new();

        for file in staged_files {
            if let Some(p) = &pattern && !p.matches_path(&file) {
                continue;
            }

            let diff_lines = ctx.staged_diff(&file)?;

            for line in diff_lines {
                if let DiffLine::Added(content) = line {
                    if regex.is_match(&content) {
                        matched_files.push(file.display().to_string());
                        break;
                    }
                }
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
    use crate::context::{DiffLine, MockContext};
    use crate::tmpl;

    #[test]
    fn test_suppress_string_success() -> Result<()> {
        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), None);
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_staged_diff()
            .returning(|_| Ok(vec![DiffLine::Added("clean content".to_string())]));

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }

    #[test]
    fn test_suppress_string_failure() -> Result<()> {
        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), None);
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_staged_diff()
            .returning(|_| Ok(vec![DiffLine::Added("this has a TODO item".to_string())]));

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
        let rule = SuppressString::new("test".to_string(), tmpl!("TODO"), Some(tmpl!("*.rs")));
        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }
}

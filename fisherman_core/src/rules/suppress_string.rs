use crate::context::{Context, DiffLine};
use crate::rules::{Rule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use regex::Regex;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct SuppressStringRule {
    pub regex: TemplateString,
    pub glob: Option<TemplateString>,
}

impl std::fmt::Display for SuppressStringRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Suppress string matching: {}", self.regex)
    }
}

#[typetag::serde(name = "suppress-string")]
impl Rule for SuppressStringRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables_new()?;
        let regex = Regex::new(&self.regex.compile(&variables)?)?;

        let pattern = match &self.glob {
            Some(g) => {
                let glob = g.compile(&variables)?;
                Some(Pattern::new(glob.as_str())?)
            }
            None => None,
        };

        let staged_files = ctx.staged_files()?;

        let mut matched_files = Vec::new();

        for file in staged_files {
            if let Some(p) = &pattern
                && !p.matches_path(&file)
            {
                continue;
            }

            let diff_lines = ctx.staged_diff(&file)?;

            for line in diff_lines {
                if let DiffLine::Added(content) = line
                    && regex.is_match(&content)
                {
                    matched_files.push(file.display().to_string());
                    break;
                }
            }
        }

        if !matched_files.is_empty() {
            return Ok(RuleResult::Failure {
                name: "suppress-string".into(),
                message: format!(
                    "The following files contain suppressed string: {}",
                    matched_files.join(", ")
                ),
            });
        }

        Ok(RuleResult::Success {
            name: "suppress-string".into(),
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
    fn serialize_test() -> Result<()> {
        let config = SuppressStringRule {
            regex: "TODO".into(),
            glob: None,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"regex":"TODO","glob":null}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: SuppressStringRule = serde_json::from_str(r#"{"regex":"TODO"}"#)?;

        assert_eq!(config.regex, "TODO".into());
        assert!(config.glob.is_none());

        Ok(())
    }

    #[test]
    fn serialize_test_with_glob() -> Result<()> {
        let config = SuppressStringRule {
            regex: "TODO".into(),
            glob: Some("*.rs".into()),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"regex":"TODO","glob":"*.rs"}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test_with_glob() -> Result<()> {
        let config: SuppressStringRule =
            serde_json::from_str(r#"{"regex":"TODO","glob":"*.rs"}"#)?;

        assert_eq!(config.regex, "TODO".into());
        assert_eq!(config.glob, Some("*.rs".into()));

        Ok(())
    }


    #[test]
    fn test_suppress_string_success() -> Result<()> {
        let rule = SuppressStringRule {
            regex: "TODO".into(),
            glob: None,
        };

        let mut context = MockContext::new();
        context
            .expect_variables_new()
            .returning(|| Ok(std::collections::HashMap::new()));
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
        let rule = SuppressStringRule {
            regex: tmpl!("TODO"),
            glob: None,
        };

        let mut context = MockContext::new();
        context
            .expect_variables_new()
            .returning(|| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));
        context
            .expect_staged_diff()
            .returning(|_| Ok(vec![DiffLine::Added("this has a TODO item".to_string())]));

        let result = rule.check(&context)?;
        match result {
            RuleResult::Failure { name, message } => {
                assert_eq!(name, "suppress-string");
                assert!(
                    message.contains("The following files contain suppressed string: test.txt")
                );
            }
            _ => panic!("Expected failure"),
        }
        Ok(())
    }

    #[test]
    fn test_suppress_string_with_glob() -> Result<()> {
        let rule = SuppressStringRule {
            regex: tmpl!("TODO"),
            glob: Some(tmpl!("*.rs")),
        };

        let mut context = MockContext::new();
        context
            .expect_variables_new()
            .returning(|| Ok(std::collections::HashMap::new()));
        context
            .expect_staged_files()
            .returning(|| Ok(vec![std::path::PathBuf::from("test.txt")]));

        let result = rule.check(&context)?;
        assert!(matches!(result, RuleResult::Success { .. }));
        Ok(())
    }

    #[test]
    fn test_suppress_string_rule_success() -> Result<()> {
        let rule = SuppressStringRule {
            regex: "TODO".into(),
            glob: None,
        };

        let mut context = MockContext::new();
        context
            .expect_variables_new()
            .returning(|| Ok(std::collections::HashMap::new()));
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
    fn test_display() {
        let rule = SuppressStringRule {
            regex: "TODO".into(),
            glob: None,
        };
        assert_eq!(format!("{}", rule), "Suppress string matching: `TODO`");
    }
}

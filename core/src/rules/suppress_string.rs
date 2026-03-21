use crate::context::{Context, DiffLine};
use crate::extract_vars;
use crate::rules::rule::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use glob::Pattern;
use regex::Regex;
use rules_derive::ConditionalRule as ConditionalRuleDerive;

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct SuppressStringRule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
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
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: "suppress-string".into(),
            });
        }

        let variables = extract_vars!(&self, ctx)?;
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
                    && regex.is_match(&content) {
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
    fn test_suppress_string_success() -> Result<()> {
        let rule = SuppressStringRule {
            when: None,
            extract: None,
            regex: "TODO".into(),
            glob: None,
        };

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
        let rule = SuppressStringRule {
            regex: tmpl!("TODO"),
            glob: None,
            when: None,
            extract: None,
        };

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
            when: None,
            extract: None,
        };

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

    #[test]
    fn test_suppress_string_rule_success() -> Result<()> {
        let rule = SuppressStringRule {
            when: None,
            extract: None,
            regex: "TODO".into(),
            glob: None,
        };

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
    fn test_display() {
        let rule = SuppressStringRule {
            when: None,
            extract: None,
            regex: "TODO".into(),
            glob: None,
        };
        assert_eq!(format!("{}", rule), "Suppress string matching: `TODO`");
    }
}

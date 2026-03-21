use crate::context::Context;
use crate::extract_vars;
use crate::rules::rule::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use rules_derive::ConditionalRule as ConditionalRuleDerive;
use std::fs::OpenOptions;
use std::io::Write;

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct WriteFileRule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    pub path: TemplateString,
    pub content: TemplateString,
    pub append: Option<bool>,
}

#[typetag::serde(name = "write-file")]
impl Rule for WriteFileRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: "write-file".into(),
            });
        }

        let variables = extract_vars!(self, ctx)?;
        let path = self.path.compile(&variables)?;
        let content = self.content.compile(&variables)?;

        let append = self.append.unwrap_or(false);

        let mut file = OpenOptions::new()
            .write(true)
            .create(true)
            .append(append)
            .open(path)?;

        file.write_all(content.as_bytes())?;

        Ok(RuleResult::Success {
            name: "write-file".into(),
            output: None,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use crate::t;
    use anyhow::Result;
    use std::collections::HashMap;
    use std::fs;
    use tempfile::TempDir;

    fn mock_ctx_with_vars(vars: HashMap<String, String>) -> MockContext {
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));
        ctx
    }

    fn mock_ctx() -> MockContext {
        mock_ctx_with_vars(HashMap::new())
    }

    #[test]
    fn write_file_when_file_doesnt_exist() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx())?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path)?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[test]
    fn write_file_when_file_exists() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test")?;

        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(false),
            when: None,
            extract: None,
        };

        let result = rule.check(&mock_ctx())?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path)?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[test]
    fn append_file_when_file_exists() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test")?;

        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(true),
            when: None,
            extract: None,
        };

        let result = rule.check(&mock_ctx())?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path)?;
        assert_eq!(file_content, "TestHello, world!");

        Ok(())
    }

    #[test]
    fn write_file_when_path_template_literal() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("{{FILE_NAME}}.txt");
        let content = "Hello, world!".to_string();

        let mut variables = HashMap::new();
        variables.insert("FILE_NAME".to_string(), "test".to_string());

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx_with_vars(variables))?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(dir.path().join("test.txt"))?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[test]
    fn write_file_when_content_template_literal() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, {{WHO}}!".to_string();

        let mut variables = HashMap::new();
        variables.insert("WHO".to_string(), "world".to_string());

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx_with_vars(variables))?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path)?;
        assert_eq!(file_content, "Hello, world!");

        Ok(())
    }

    #[test]
    fn test_write_file_path_template_error() {
        let rule = WriteFileRule {
            path: t!("{{missing}}/file.txt"),
            content: t!("content"),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx());
        assert!(result.is_err());
    }

    #[test]
    fn test_write_file_content_template_error() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!("{{missing}}"),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx());
        assert!(result.is_err());

        Ok(())
    }

    #[test]
    fn test_write_file_io_error() {
        let rule = WriteFileRule {
            path: t!("/invalid/path/that/does/not/exist/file.txt"),
            content: t!("content"),
            when: None,
            extract: None,
            append: Some(false),
        };

        let result = rule.check(&mock_ctx());
        assert!(result.is_err());
    }

    #[test]
    fn test_write_file_rule_new() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            when: None,
            extract: None,
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: None,
        };

        let result = rule.check(&mock_ctx())?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path)?;
        assert_eq!(file_content, content);

        Ok(())
    }
}

use crate::context::Context;
use crate::rules::{ExecutionMode, Rule, RuleResult};
use anyhow::Result;
use async_trait::async_trait;
use template_str::TemplateString;
use tokio::fs::OpenOptions;
use tokio::io::AsyncWriteExt;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct WriteFileRule {
    pub path: TemplateString,
    pub content: TemplateString,
    pub append: Option<bool>,
}

impl std::fmt::Display for WriteFileRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Write file: {}", self.path)
    }
}

#[async_trait]
#[typetag::serde(name = "write-file")]
impl Rule for WriteFileRule {
    fn execution_mode(&self) -> ExecutionMode {
        ExecutionMode::Async
    }

    async fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables()?;
        let path = self.path.compile(&variables)?;
        let content = self.content.compile(&variables)?;

        let append = self.append.unwrap_or(false);

        let mut file = OpenOptions::new()
            .write(true)
            .create(true)
            .append(append)
            .open(path)
            .await?;

        file.write_all(content.as_bytes()).await?;
        file.flush().await?;

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
    use anyhow::Result;
    use std::collections::HashMap;
    use tempfile::TempDir;
    use template_str::t;
    use tokio::fs;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("hello"),
            append: Some(false),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"path":"/tmp/output.txt","content":"hello","append":false}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: WriteFileRule =
            serde_json::from_str(r#"{"path":"/tmp/output.txt","content":"hello","append":true}"#)?;

        assert_eq!(config.path, t!("/tmp/output.txt"));
        assert_eq!(config.content, t!("hello"));
        assert_eq!(config.append, Some(true));

        Ok(())
    }

    #[test]
    fn serialize_test_with_append_true() -> Result<()> {
        let config = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("hello"),
            append: Some(true),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"path":"/tmp/output.txt","content":"hello","append":true}"#
        );

        Ok(())
    }

    #[test]
    fn serialize_test_with_append_none() -> Result<()> {
        let config = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("hello"),
            append: None,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"path":"/tmp/output.txt","content":"hello","append":null}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test_with_append_none() -> Result<()> {
        let config: WriteFileRule =
            serde_json::from_str(r#"{"path":"/tmp/output.txt","content":"hello"}"#)?;

        assert_eq!(config.path, t!("/tmp/output.txt"));
        assert_eq!(config.content, t!("hello"));
        assert!(config.append.is_none());

        Ok(())
    }

    #[test]
    fn serialize_test_with_append() -> Result<()> {
        let config = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("hello"),
            append: Some(false),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"path":"/tmp/output.txt","content":"hello","append":false}"#
        );

        Ok(())
    }

    #[test]
    fn serialize_test_with_all_fields() -> Result<()> {
        let config = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("hello"),
            append: Some(true),
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"path":"/tmp/output.txt","content":"hello","append":true}"#
        );

        Ok(())
    }

    fn mock_ctx_with_vars(vars: HashMap<String, String>) -> MockContext {
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move || Ok(vars.clone()));
        ctx
    }

    fn mock_ctx() -> MockContext {
        mock_ctx_with_vars(HashMap::new())
    }

    #[tokio::test]
    async fn write_file_when_file_doesnt_exist() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx()).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).await?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[tokio::test]
    async fn write_file_when_file_exists() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test").await?;

        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx()).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).await?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[tokio::test]
    async fn append_file_when_file_exists() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test").await?;

        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(true),
        };

        let result = rule.check(&mock_ctx()).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).await?;
        assert_eq!(file_content, "TestHello, world!");

        Ok(())
    }

    #[tokio::test]
    async fn write_file_when_path_template_literal() -> Result<()> {
        let dir = TempDir::new()?;

        let path = dir.path().join("{{FILE_NAME}}.txt");
        let content = "Hello, world!".to_string();

        let mut variables = HashMap::new();
        variables.insert("FILE_NAME".to_string(), "test".to_string());

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx_with_vars(variables)).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(dir.path().join("test.txt")).await?;
        assert_eq!(file_content, content);

        Ok(())
    }

    #[tokio::test]
    async fn write_file_when_content_template_literal() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, {{WHO}}!".to_string();

        let mut variables = HashMap::new();
        variables.insert("WHO".to_string(), "world".to_string());

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx_with_vars(variables)).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).await?;
        assert_eq!(file_content, "Hello, world!");

        Ok(())
    }

    #[tokio::test]
    async fn test_write_file_path_template_error() {
        let rule = WriteFileRule {
            path: t!("{{missing}}/file.txt"),
            content: t!("content"),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx()).await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_write_file_content_template_error() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!("{{missing}}"),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx()).await;
        assert!(result.is_err());

        Ok(())
    }

    #[tokio::test]
    async fn test_write_file_io_error() {
        let rule = WriteFileRule {
            path: t!("/invalid/path/that/does/not/exist/file.txt"),
            content: t!("content"),
            append: Some(false),
        };

        let result = rule.check(&mock_ctx()).await;
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = WriteFileRule {
            path: t!("/tmp/output.txt"),
            content: t!("content"),
            append: None,
        };
        assert_eq!(format!("{}", rule), "Write file: `/tmp/output.txt`");
    }

    #[tokio::test]
    async fn test_write_file_rule_new() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFileRule {
            path: t!(path.to_str().unwrap()),
            content: t!(content.clone()),
            append: None,
        };

        let result = rule.check(&mock_ctx()).await?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write-file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).await?;
        assert_eq!(file_content, content);

        Ok(())
    }
}

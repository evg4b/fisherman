use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use std::fs::OpenOptions;
use std::io::Write;

pub struct WriteFile {
    name: String,
    path: TemplateString,
    content: TemplateString,
    append: bool,
}

impl WriteFile {
    pub fn new(
        name: String,
        path: TemplateString,
        content: TemplateString,
        append: bool,
    ) -> WriteFile {
        WriteFile {
            name,
            path,
            content,
            append,
        }
    }
}

impl CompiledRule for WriteFile {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let variables = ctx.variables(&[])?;
        let content = self.content.to_string(&variables)?;
        let path = self.path.to_string(&variables)?;

        let mut file = OpenOptions::new()
            .write(true)
            .create(true)
            .append(self.append)
            .open(&path)?;

        file.write_all(content.as_bytes())?;

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
    use crate::t;
    use std::collections::HashMap;
    use std::fs;
    use anyhow::Result;
    use tempfile::TempDir;

    #[test]
    fn write_file_when_file_doesnt_exist() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write_file");
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

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write_file");
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

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            true,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write_file");
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

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| {
            let mut variables = HashMap::new();
            variables.insert("FILE_NAME".to_string(), "test".to_string());
            Ok(variables)
        });

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write_file");
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

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, None);

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, "Hello, world!");

        Ok(())
    }

    #[test]
    fn test_sync() {
        let rule = WriteFile::new(
            "write_file".to_string(),
            t!("path/to/file.txt"),
            t!("content"),
            false,
        );
        assert!(!rule.sync());
    }

    #[test]
    fn test_write_file_variables_error() {
        let rule = WriteFile::new(
            "write_file".to_string(),
            t!("path/to/file.txt"),
            t!("content"),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_write_file_path_template_error() {
        let rule = WriteFile::new(
            "write_file".to_string(),
            t!("{{missing}}/file.txt"),
            t!("content"),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_write_file_content_template_error() {
        let dir = TempDir::new("write_file_template_error").unwrap();
        let path = dir.path().join("test.txt");

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!("{{missing}}"),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }

    #[test]
    fn test_write_file_io_error() {
        let rule = WriteFile::new(
            "write_file".to_string(),
            t!("/invalid/path/that/does/not/exist/file.txt"),
            t!("content"),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&context);
        assert!(result.is_err());
    }
}

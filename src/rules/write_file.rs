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
    fn check(&self, _: &dyn Context) -> Result<RuleResult> {
        let content = self.content.to_string()?;
        let path = self.path.to_string()?;

        let mut file = OpenOptions::new()
            .write(true)
            .create(true)
            .append(self.append)
            .open(&path)?;

        file.write_all(content.as_bytes())?;

        Ok(RuleResult::Success {
            name: self.name.clone(),
            output: "".to_string(),
        })
    }
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;
    use super::*;
    use crate::context::MockContext;
    use std::fs;
    use anyhow::Result;
    use crate::tmpl;
    use tempfile::TempDir;

    #[test]
    fn write_file_when_file_doesnt_exist() -> Result<()> {
        let dir = TempDir::new()?;
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();
        let variables = HashMap::new();

        let rule = WriteFile::new(
            "write_file".to_string(),
            tmpl!(path.to_str().as_ref().unwrap(), variables.clone()),
            tmpl!(content, variables),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

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
        let variables = HashMap::new();

        let rule = WriteFile::new(
            "write_file".to_string(),
            tmpl!(path.to_str().as_ref().unwrap(), variables.clone()),
            tmpl!(content.clone(), variables),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

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
        let variables = HashMap::new();

        let rule = WriteFile::new(
            "write_file".to_string(),
            tmpl!(path.to_str().as_ref().unwrap(), variables.clone()),
            tmpl!(content.clone(), variables),
            true,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

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

        let rule = WriteFile::new(
            "write_file".to_string(),
            tmpl!(path.to_str().as_ref().unwrap(), variables.clone()),
            tmpl!(content.clone(), variables),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

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

        let rule = WriteFile::new(
            "write_file".to_string(),
            tmpl!(path.to_str().as_ref().unwrap(), variables.clone()),
            tmpl!(content.clone(), variables),
            false,
        );

        let context = MockContext::new();
        let result = rule.check(&context)?;

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, "Hello, world!");

        Ok(())
    }
}

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
            output: "".to_string(),
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use std::collections::HashMap;
    use std::fs;
    use tempdir::TempDir;
    use crate::t;

    #[test]
    fn write_file_when_file_doesnt_exist() {
        let dir = TempDir::new("write_file_when_file_doesnt_exist").unwrap();
        let path = dir.path().join("test.txt");
        let content = "Hello, world!".to_string();

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&context).unwrap();

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, content);
    }

    #[test]
    fn write_file_when_file_exists() {
        let dir = TempDir::new("write_file_when_file_exists").unwrap();

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test").unwrap();

        let content = "Hello, world!".to_string();

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.check(&context).unwrap();

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, content);
    }

    #[test]
    fn append_file_when_file_exists() {
        let dir = TempDir::new("write_file_when_file_exists").unwrap();

        let path = dir.path().join("test.txt");
        fs::write(&path, "Test").unwrap();

        let content = "Hello, world!".to_string();

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            true,
        );

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));
        let result = rule.check(&context).unwrap();

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, "TestHello, world!");
    }

    #[test]
    fn write_file_when_path_template_literal() {
        let dir = TempDir::new("write_file_when_file_doesnt_exist").unwrap();

        let path = dir.path().join("{{FILE_NAME}}.txt");
        let content = "Hello, world!".to_string();

        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| {
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

        let result = rule.check(&context).unwrap();

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(dir.path().join("test.txt")).unwrap();
        assert_eq!(file_content, content);
    }

    #[test]
    fn write_file_when_content_template_literal() {
        let dir = TempDir::new("write_file_when_file_doesnt_exist").unwrap();
        let path = dir.path().join("test.txt");
        let content = "Hello, {{WHO}}!".to_string();

        let rule = WriteFile::new(
            "write_file".to_string(),
            t!(path.to_str().unwrap()),
            t!(content.clone()),
            false,
        );


        let mut context = MockContext::new();
        context
            .expect_variables()
            .returning(|_| {
                let mut variables = HashMap::new();
                variables.insert("WHO".to_string(), "world".to_string());
                Ok(variables)
            });
        let result = rule.check(&context).unwrap();

        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };
        assert_eq!(name, "write_file");
        assert_eq!(output, "");

        let file_content = fs::read_to_string(path).unwrap();
        assert_eq!(file_content, "Hello, world!");
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
}

use std::collections::HashMap;
use crate::templates::{replace_in_string, TemplateError};

#[derive(Debug)]
pub struct TemplateString<'a> {
    pub template: String,
    pub variables: &'a HashMap<String, String>,
}

impl TemplateString<'_> {
    pub fn new(template: String, variables: &HashMap<String, String>) -> TemplateString {
        TemplateString { template, variables }
    }

    pub fn to_string(&self) -> Result<String, TemplateError> {
        replace_in_string(&self.template, self.variables)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    #[test]
    fn test_template_string_with_placeholders() {
        let template = "Hello, {{name}}! Welcome to {{place}}.".to_string();
        let mut variables = HashMap::new();
        variables.insert("name".to_string(), "Alice".to_string());
        variables.insert("place".to_string(), "Wonderland".to_string());

        let template_string = TemplateString::new(template, &variables);
        let result = template_string.to_string().unwrap();

        assert_eq!(result, "Hello, Alice! Welcome to Wonderland.");
    }

    #[test]
    fn test_template_string_without_placeholders() {
        let template = "Hello, World!".to_string();
        let variables = HashMap::new();

        let template_string = TemplateString::new(template, &variables);
        let result = template_string.to_string().unwrap();

        assert_eq!(result, "Hello, World!");
    }

    #[test]
    fn test_template_string_with_missing_placeholder() {
        let template = "Hello, {{name}}! Welcome to {{place}}.".to_string();
        let mut variables = HashMap::new();
        variables.insert("name".to_string(), "Alice".to_string());
        // Missing "place" variable

        let template_string = TemplateString::new(template, &variables);
        let result = template_string.to_string();

        assert!(result.is_err());
        match result {
            Err(TemplateError::PlaceholderNotFound { placeholder }) => {
                assert_eq!(placeholder, "place");
            }
            _ => panic!("Expected PlaceholderNotFound error"),
        }
    }
}
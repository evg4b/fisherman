use crate::templates::TemplateError;
use std::collections::HashMap;

#[derive(Debug)]
pub struct TemplateString {
    template: String,
}

impl TemplateString {
    pub fn from<T>(template: T) -> Self
    where
        T: Into<String>,
    {
        TemplateString {
            template: template.into(),
        }
    }

    pub fn to_string(&self, variables: &HashMap<String, String>) -> Result<String, TemplateError> {
        let input = self.template.as_ref();
        let mut result = self.template.to_owned();
        let pattern = regex::Regex::new(r"\{\{(.*?)}}").unwrap();

        for cap in pattern.captures_iter(input) {
            if let Some(key) = cap.get(1) {
                let key_str = key.as_str();
                match variables.get(key_str) {
                    Some(value) => {
                        result = result.replace(&format!("{{{{{}}}}}", key_str), value);
                    }
                    None => {
                        return Err(TemplateError::PlaceholderNotFound {
                            placeholder: key_str.to_string(),
                        });
                    }
                }
            }
        }
        Ok(result)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_template_string_creation() {
        let template = t!("test");

        assert_eq!(template.template, "test");
    }

    #[test]
    fn test_template_string_no_placeholders() {
        let variables = HashMap::new();
        let template = t!("Hello, world!");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "Hello, world!");
    }

    #[test]
    fn test_template_string_with_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));
        variables.insert(String::from("greeting"), String::from("Hello"));

        let template = t!("{{greeting}}, {{name}}!");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "Hello, John!");
    }

    #[test]
    fn test_template_string_multiple_same_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));

        let template = t!("Hello, {{name}}! How are you, {{name}}?");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "Hello, John! How are you, John?");
    }

    #[test]
    fn test_template_string_missing_placeholder() {
        let mut variables = HashMap::new();
        variables.insert(String::from("greeting"), String::from("Hello"));

        let template = t!("{{greeting}}, {{name}}!");

        let result = template.to_string(&variables);
        assert!(result.is_err());

        match result.unwrap_err() {
            TemplateError::PlaceholderNotFound { placeholder } => {
                assert_eq!(placeholder, "name");
            }
            _ => panic!("Expected PlaceholderNotFound error"),
        }
    }

    #[test]
    fn test_template_string_empty_template() {
        let variables = HashMap::new();
        let template = t!("");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "");
    }

    #[test]
    fn test_template_string_with_partial_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));

        let template = t!("Hello, {{name}}!");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "Hello, John!");
    }

    #[test]
    fn test_template_string_with_braces_in_value() {
        let mut variables = HashMap::new();
        variables.insert(String::from("value"), String::from("a {nested} value"));

        let template = t!("This is {{value}}");

        let result = template.to_string(&variables).unwrap();
        assert_eq!(result, "This is a {nested} value");
    }
}

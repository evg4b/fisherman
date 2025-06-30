use crate::templates::TemplateError;
use std::collections::HashMap;

#[derive(Debug)]
pub struct TemplateStringLegacy {
    template: String,
    // TODO: Reduce memory consumption by using a link to the map instead of a copy
    variables: HashMap<String, String>,
}

impl TemplateStringLegacy {
    pub fn new(template: String, variables: HashMap<String, String>) -> TemplateStringLegacy {
        TemplateStringLegacy { template, variables }
    }

    pub fn to_string(&self) -> Result<String, TemplateError> {
        let input = self.template.as_ref();
        let mut result = self.template.clone();
        let pattern = regex::Regex::new(r"\{\{(.*?)}}").unwrap();

        for cap in pattern.captures_iter(input) {
            if let Some(key) = cap.get(1) {
                let key_str = key.as_str();
                match self.variables.get(key_str) {
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
    use crate::tmpl_legacy;
    use std::collections::HashMap;

    #[test]
    fn test_template_string_creation() {
        let variables = HashMap::new();
        let template = tmpl_legacy!(String::from("test"), variables);

        assert_eq!(template.template, "test");
        assert_eq!(template.variables.len(), 0);
    }

    #[test]
    fn test_template_string_no_placeholders() {
        let variables = HashMap::new();
        let template = tmpl_legacy!(String::from("Hello, world!"), variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "Hello, world!");
    }

    #[test]
    fn test_template_string_with_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));
        variables.insert(String::from("greeting"), String::from("Hello"));

        let template = tmpl_legacy!("{{greeting}}, {{name}}!", variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "Hello, John!");
    }

    #[test]
    fn test_template_string_multiple_same_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));

        let template = tmpl_legacy!(String::from("Hello, {{name}}! How are you, {{name}}?"), variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "Hello, John! How are you, John?");
    }

    #[test]
    fn test_template_string_missing_placeholder() {
        let mut variables = HashMap::new();
        variables.insert(String::from("greeting"), String::from("Hello"));

        let template = tmpl_legacy!(String::from("{{greeting}}, {{name}}!"), variables);

        let result = template.to_string();
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
        let template = tmpl_legacy!(String::from(""), variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "");
    }

    #[test]
    fn test_template_string_with_partial_placeholders() {
        let mut variables = HashMap::new();
        variables.insert(String::from("name"), String::from("John"));

        let template = tmpl_legacy!(String::from("Hello, {{name}}!"), variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "Hello, John!");
    }

    #[test]
    fn test_template_string_with_braces_in_value() {
        let mut variables = HashMap::new();
        variables.insert(String::from("value"), String::from("a {nested} value"));

        let template = tmpl_legacy!(String::from("This is {{value}}"), variables);

        let result = template.to_string().unwrap();
        assert_eq!(result, "This is a {nested} value");
    }
}
use crate::templates::errors::TemplateError;
use std::collections::HashMap;

pub fn replace_in_string<T: AsRef<str>>(
    input: T,
    values: &HashMap<String, String>,
) -> Result<String, TemplateError> {
    let input = input.as_ref();
    let mut result = input.to_string();
    let pattern = regex::Regex::new(r"\{\{(.*?)}}").unwrap();

    for cap in pattern.captures_iter(input) {
        if let Some(key) = cap.get(1) {
            let key_str = key.as_str();
            match values.get(key_str) {
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

#[cfg(test)]
mod test_replace_placeholders {
    use super::*;

    #[test]
    fn should_replace_placeholders_in_string() {
        let mut values = HashMap::new();
        values.insert("name".to_string(), "World".to_string());
        values.insert("greeting".to_string(), "Hello".to_string());

        let template = "{{greeting}}, {{name}}!";
        let result = replace_in_string(template, &values).unwrap();
        assert_eq!(result, "Hello, World!");
    }

    #[test]
    fn should_return_error_if_key_not_found() {
        let mut values = HashMap::new();
        values.insert("greeting".to_string(), "Hello".to_string());
        let template = "{{greeting}}, {{name}}!";
        let result = replace_in_string(template, &values);
        assert!(result.is_err());
    }
}

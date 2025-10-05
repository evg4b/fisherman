use crate::t;
use crate::templates::TemplateError;
use std::collections::HashMap;

pub fn replace_in_vac(
    input: &[String],
    values: &HashMap<String, String>,
) -> Result<Vec<String>, TemplateError> {
    let transformed: Vec<String> = input
        .iter()
        .map(|v| t!(v).to_string(values))
        .collect::<Result<Vec<String>, TemplateError>>()?;

    Ok(transformed)
}

#[cfg(test)]
mod replace_in_vac_test {
    use super::*;
    use std::collections::HashMap;

    #[test]
    fn should_replace_placeholders_in_vac() {
        let mut values = HashMap::new();
        values.insert("name".to_string(), "World".to_string());
        values.insert("greeting".to_string(), "Hello".to_string());

        let input = vec!["{{name}}".to_string(), "{{greeting}}".to_string()];

        let result = replace_in_vac(&input, &values).unwrap();
        assert_eq!(result[0], "World");
        assert_eq!(result[1], "Hello");
    }

    #[test]
    fn should_return_error_if_key_not_found() {
        let mut values = HashMap::new();
        values.insert("greeting".to_string(), "Hello".to_string());

        let input = vec!["{{name}}".to_string(), "{{greeting}}".to_string()];

        let result = replace_in_vac(&input, &values);
        assert!(result.is_err());
    }
}

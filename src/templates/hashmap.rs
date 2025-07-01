use crate::t;
use crate::templates::TemplateError;
use std::collections::HashMap;

pub fn replace_in_hashmap(
    input: &HashMap<String, String>,
    values: &HashMap<String, String>,
) -> Result<HashMap<String, String>, TemplateError> {
    let transformed: HashMap<String, String> = input
        .iter()
        .map(|(k, v)| match t!(v).to_string(values) {
            Ok(v) => Ok((k.to_owned(), v)),
            Err(e) => match e {
                TemplateError::PlaceholderNotFound { placeholder } => {
                    Err(TemplateError::PlaceholderNotFoundForKey {
                        placeholder,
                        key: k.clone(),
                    })
                }
                _ => Err(e),
            },
        })
        .collect::<Result<HashMap<String, String>, TemplateError>>()?;

    Ok(transformed)
}

#[cfg(test)]
mod test_replace_placeholders_in_hashmap {
    use super::*;
    use std::collections::HashMap;

    #[test]
    fn should_replace_placeholders_in_hashmap() {
        let mut values = HashMap::new();
        values.insert("name".to_string(), "World".to_string());
        values.insert("greeting".to_string(), "Hello".to_string());

        let mut input = HashMap::new();
        input.insert("name".to_string(), "{{name}}".to_string());
        input.insert("greeting".to_string(), "{{greeting}}".to_string());

        let result = replace_in_hashmap(&input, &values).unwrap();
        assert_eq!(result.get("name").unwrap(), "World");
        assert_eq!(result.get("greeting").unwrap(), "Hello");
    }

    #[test]
    fn should_return_error_if_key_not_found() {
        let mut values = HashMap::new();
        values.insert("greeting".to_string(), "Hello".to_string());

        let mut input = HashMap::new();
        input.insert("name".to_string(), "{{name}}".to_string());
        input.insert("greeting".to_string(), "{{greeting}}".to_string());

        let result = replace_in_hashmap(&input, &values);
        assert!(result.is_err());
    }
}

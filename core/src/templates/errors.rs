use std::fmt::{Display, Formatter};

#[derive(Debug, Clone, PartialEq)]
pub enum TemplateError {
    PlaceholderNotFound { placeholder: String },
    PlaceholderNotFoundForKey { placeholder: String, key: String },
}

impl Display for TemplateError {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            TemplateError::PlaceholderNotFound { placeholder } => {
                write!(f, "Variable '{}' not defined", placeholder)
            }
            TemplateError::PlaceholderNotFoundForKey { placeholder, key } => {
                write!(f, "Variable '{}' not found for key '{}'", placeholder, key)
            }
        }
    }
}

impl std::error::Error for TemplateError {}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_placeholder_not_found_display() {
        let error = TemplateError::PlaceholderNotFound {
            placeholder: "my_var".to_string(),
        };
        assert_eq!(error.to_string(), "Variable 'my_var' not defined");
    }

    #[test]
    fn test_placeholder_not_found_for_key_display() {
        let error = TemplateError::PlaceholderNotFoundForKey {
            placeholder: "my_var".to_string(),
            key: "env_key".to_string(),
        };
        assert_eq!(
            error.to_string(),
            "Variable 'my_var' not found for key 'env_key'"
        );
    }
}

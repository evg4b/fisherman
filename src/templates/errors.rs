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
                write!(f, "Placeholder '{}' not found", placeholder)
            }
            TemplateError::PlaceholderNotFoundForKey { placeholder, key } => {
                write!(f, "Placeholder '{}' not found for key '{}'", placeholder, key)
            }
        }
    }
}

impl std::error::Error for TemplateError {}
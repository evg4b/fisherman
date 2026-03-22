use std::fmt;
use std::path::PathBuf;

#[derive(Debug)]
pub(crate) enum HookError {
    AlreadyExists { name: &'static str, hook: PathBuf },
}

impl fmt::Display for HookError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            HookError::AlreadyExists { hook, name } => {
                write!(f, "Hook '{}' already installed at {}", name, hook.display())
            }
        }
    }
}

impl std::error::Error for HookError {}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::PathBuf;

    #[test]
    fn test_already_exists_display() {
        let error = HookError::AlreadyExists {
            name: "pre-commit",
            hook: PathBuf::from("/tmp/.git/hooks/pre-commit"),
        };
        let msg = error.to_string();
        assert!(msg.contains("pre-commit"));
        assert!(msg.contains("/tmp/.git/hooks/pre-commit"));
    }
}

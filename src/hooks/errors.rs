use std::fmt;
use std::path::PathBuf;

#[derive(Debug)]
pub(crate) enum HookError {
    AlreadyExists { hook: PathBuf },
}

impl fmt::Display for HookError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            HookError::AlreadyExists { hook } => {
                write!(f, "Hook {} already exists", hook.display())
            }
        }
    }
}

impl std::error::Error for HookError {}

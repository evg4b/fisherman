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
                write!(f, "Hook {} already exists ({})", name, hook.display())
            }
        }
    }
}

impl std::error::Error for HookError {}

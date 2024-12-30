use std::fmt;
use std::path::PathBuf;

#[derive(Debug)]
pub(crate) enum ConfigurationError {
    MultipleConfigFiles { files: Vec<PathBuf> },
    UnknownConfigFileExtension,
}

impl fmt::Display for ConfigurationError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConfigurationError::MultipleConfigFiles { files } => {
                write!(f, "Multiple configuration files found:\n")?;
                files.iter().for_each(|file| {
                    write!(f, "  {}\n", file.display()).unwrap();
                });
                Ok(())
            }
            ConfigurationError::UnknownConfigFileExtension => {
                write!(f, "Unknown configuration file extension")
            }
        }
    }
}

impl std::error::Error for ConfigurationError {}

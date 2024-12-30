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
                writeln!(f, "Multiple configuration files found:")?;
                files.iter().for_each(|file| {
                    writeln!(f, "  {:?}", file).unwrap();
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

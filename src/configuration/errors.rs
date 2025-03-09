use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigurationError {
    #[error("Multiple configuration files found: {files:?}")]
    MultipleConfigFiles { files: Vec<PathBuf> },
    #[error("Unknown configuration file extension")]
    UnknownConfigFileExtension,
}

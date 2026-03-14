use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigurationError {
    #[error("Multiple config files found: {files:?}. Keep only one.")]
    MultipleConfigFiles { files: Vec<PathBuf> },
    #[error("Unsupported config file format. Use .toml, .yaml, .yml, or .json")]
    UnknownConfigFileExtension,
}

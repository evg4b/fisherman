use std::path::PathBuf;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ConfigurationError {
    #[error("Multiple config files found: {files:?}. Keep only one.")]
    MultipleConfigFiles { files: Vec<PathBuf> },
    #[error("Unsupported config file format. Use .toml, .yaml, .yml, or .json")]
    UnknownConfigFileExtension,
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::PathBuf;

    #[test]
    fn test_multiple_config_files_display() {
        let error = ConfigurationError::MultipleConfigFiles {
            files: vec![
                PathBuf::from(".fisherman.toml"),
                PathBuf::from(".fisherman.yaml"),
            ],
        };
        let msg = error.to_string();
        assert!(msg.contains("Multiple config files found"));
        assert!(msg.contains(".fisherman.toml"));
        assert!(msg.contains(".fisherman.yaml"));
    }

    #[test]
    fn test_unknown_config_file_extension_display() {
        let error = ConfigurationError::UnknownConfigFileExtension;
        assert_eq!(
            error.to_string(),
            "Unsupported config file format. Use .toml, .yaml, .yml, or .json"
        );
    }
}

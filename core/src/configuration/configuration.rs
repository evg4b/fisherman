use crate::configuration::errors::ConfigurationError;
use crate::configuration::files::find_config_files;
use crate::hooks::GitHook;
use crate::rules::Rule;
use anyhow::{bail, Result};
use figment::providers::{Format, Json, Toml, Yaml};
use figment::Figment;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::ffi::OsStr;
use std::path::{Path, PathBuf};

#[derive(Debug, Default, Deserialize, Serialize)]
pub struct Configuration {
    #[serde(default)]
    pub hooks: HashMap<GitHook, Vec<Rule>>,
    #[serde(default)]
    pub extract: Vec<String>,
    #[serde(skip)]
    pub files: Vec<PathBuf>,
}

impl Configuration {
    pub(crate) fn load(path: &Path) -> Result<Configuration> {
        let files = find_config_files(path)?;

        let mut instance = Figment::new();
        for file in files.iter() {
            let extension = match file.extension().and_then(OsStr::to_str) {
                Some(ext) => ext,
                None => bail!(ConfigurationError::UnknownConfigFileExtension),
            };

            instance = match extension {
                "toml" => instance.adjoin(Toml::file(file)),
                "yaml" | "yml" => instance.adjoin(Yaml::file(file)),
                "json" => instance.adjoin(Json::file(file)),
                _ => bail!(ConfigurationError::UnknownConfigFileExtension),
            };
        }

        let mut inner_config: Configuration = instance.extract()?;
        inner_config.files = files;

        Ok(inner_config)
    }

    pub fn get_configured_hooks(&self) -> Option<Vec<GitHook>> {
        if self.hooks.is_empty() {
            return None;
        }

        Some(self.hooks.keys().cloned().collect())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tempfile::TempDir;

    #[test]
    fn test_load_empty_dir_succeeds() {
        let dir = TempDir::new().unwrap();
        let result = Configuration::load(dir.path());
        assert!(result.is_ok());
    }

    #[test]
    fn test_load_toml_config() {
        let dir = TempDir::new().unwrap();
        let config = r#"
[[hooks.pre-commit]]
type = "message-regex"
regex = "^feat"
"#;
        fs::write(dir.path().join(".fisherman.toml"), config).unwrap();
        let result = Configuration::load(dir.path()).unwrap();
        assert!(result.hooks.contains_key(&GitHook::PreCommit));
    }

    #[test]
    fn test_load_yaml_config() {
        let dir = TempDir::new().unwrap();
        let config = "hooks:\n  pre-commit:\n    - type: message-regex\n      regex: '^feat'\n";
        fs::write(dir.path().join(".fisherman.yaml"), config).unwrap();
        let result = Configuration::load(dir.path()).unwrap();
        assert!(result.hooks.contains_key(&GitHook::PreCommit));
    }

    #[test]
    fn test_load_yml_config() {
        let dir = TempDir::new().unwrap();
        let config = "hooks:\n  pre-commit:\n    - type: message-regex\n      regex: '^feat'\n";
        fs::write(dir.path().join(".fisherman.yml"), config).unwrap();
        let result = Configuration::load(dir.path()).unwrap();
        assert!(result.hooks.contains_key(&GitHook::PreCommit));
    }

    #[test]
    fn test_load_json_config() {
        let dir = TempDir::new().unwrap();
        let config =
            r#"{"hooks": {"pre-commit": [{"type": "message-regex", "regex": "^feat"}]}}"#;
        fs::write(dir.path().join(".fisherman.json"), config).unwrap();
        let result = Configuration::load(dir.path()).unwrap();
        assert!(result.hooks.contains_key(&GitHook::PreCommit));
    }

    #[test]
    fn test_load_multiple_configs_in_same_dir_errors() {
        let dir = TempDir::new().unwrap();
        fs::write(dir.path().join(".fisherman.toml"), "[hooks]\n").unwrap();
        fs::write(dir.path().join(".fisherman.yaml"), "hooks: {}\n").unwrap();
        let result = Configuration::load(dir.path());
        assert!(result.is_err());
    }

    #[test]
    fn test_load_files_are_populated() {
        let dir = TempDir::new().unwrap();
        let config_path = dir.path().join(".fisherman.toml");
        fs::write(&config_path, "[hooks]\n").unwrap();
        let result = Configuration::load(dir.path()).unwrap();
        assert!(result.files.contains(&config_path));
    }

    #[test]
    fn test_get_configured_hooks_empty() {
        let config = Configuration::default();
        assert!(config.get_configured_hooks().is_none());
    }

    #[test]
    fn test_get_configured_hooks_non_empty() {
        let mut config = Configuration::default();
        config.hooks.insert(GitHook::PreCommit, vec![]);
        let hooks = config.get_configured_hooks();
        assert!(hooks.is_some());
        assert_eq!(hooks.unwrap().len(), 1);
    }

    #[test]
    fn test_get_configured_hooks_multiple() {
        let mut config = Configuration::default();
        config.hooks.insert(GitHook::PreCommit, vec![]);
        config.hooks.insert(GitHook::CommitMsg, vec![]);
        let hooks = config.get_configured_hooks();
        assert!(hooks.is_some());
        assert_eq!(hooks.unwrap().len(), 2);
    }
}

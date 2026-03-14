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
                "yaml" => instance.adjoin(Yaml::file(file)),
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

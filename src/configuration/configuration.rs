use crate::configuration::errors::ConfigurationError;
use crate::configuration::files::find_config_files;
use crate::hooks::GitHook;
use crate::rules::Rule;
use anyhow::{bail, Result};
use figment::providers::{Format, Json, Toml, Yaml};
use figment::Figment;
use serde::Deserialize;
use std::collections::HashMap;
use std::ffi::OsStr;
use std::path::{Path, PathBuf};

#[derive(Debug, Default, Deserialize)]
struct InnerConfiguration {
    pub hooks: Option<HashMap<GitHook, Vec<Rule>>>,
    pub extract: Option<Vec<String>>,
}

pub struct Configuration {
    pub hooks: HashMap<GitHook, Vec<Rule>>,
    pub files: Vec<PathBuf>,
    pub extract: Vec<String>,
}

impl Configuration {
    pub(crate) fn load(path: &Path) -> Result<Configuration> {
        let files = find_config_files(path)?;

        let mut instance = Figment::new();
        for file in files.clone() {
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

        let inner_config: InnerConfiguration = instance.extract()?;

        Ok(Configuration {
            hooks: inner_config.hooks.unwrap_or_default(),
            extract: inner_config.extract.unwrap_or_default(),
            files,
        })
    }

    pub fn get_configured_hooks(&self) -> Option<Vec<GitHook>> {
        if self.hooks.is_empty() {
            return None;
        }
        
        Some(self.hooks.keys().cloned().collect())
    }
}

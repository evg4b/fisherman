pub(crate) mod errors;

use crate::common::BError;
use crate::configuration::errors::ConfigurationError;
use crate::err;
use crate::hooks::GitHook;
use crate::rules::RuleRef;
use dirs::home_dir;
use figment::providers::{Format, Json, Toml, Yaml};
use figment::Figment;
use serde::Deserialize;
use std::collections::HashMap;
use std::ffi::OsStr;
use std::path::PathBuf;

#[derive(Debug, Default, Deserialize)]
struct InnerConfiguration {
    pub hooks: Option<HashMap<GitHook, Vec<RuleRef>>>,
}

pub(crate) struct Configuration {
    pub hooks: HashMap<GitHook, Vec<RuleRef>>,
    pub files: Vec<PathBuf>,
}

impl Configuration {
    pub(crate) fn load(path: &PathBuf) -> Result<Configuration, BError> {
        let files = find_config_files(path.clone())?;

        let mut instance = Figment::new();
        for file in files.clone() {
            let extension = match file.extension().and_then(OsStr::to_str) {
                Some(ext) => ext,
                None => err!(ConfigurationError::UnknownConfigFileExtension),
            };

            instance = match extension {
                "toml" => instance.adjoin(Toml::file(file)),
                "yaml" => instance.adjoin(Yaml::file(file)),
                "json" => instance.adjoin(Json::file(file)),
                _ => err!(ConfigurationError::UnknownConfigFileExtension),
            };
        }

        let inner_config: InnerConfiguration = instance.extract()?;

        Ok(Configuration {
            hooks: inner_config.hooks.unwrap_or_default(),
            files,
        })
    }
}

fn find_config_files(path: PathBuf) -> Result<Vec<PathBuf>, BError> {
    let locations = vec![path.join(".git"), path.clone(), home_dir().unwrap()];
    let mut config_files = vec![];
    for location_path in locations {
        let files = resolve_configs(location_path);
        if files.len() > 1 {
            err!(ConfigurationError::MultipleConfigFiles { files });
        } else if files.len() == 1 {
            config_files.push(files[0].clone());
        }
    }

    Ok(config_files)
}

fn resolve_configs(path: PathBuf) -> Vec<PathBuf> {
    vec![
        path.join(".fisherman.toml"),
        path.join(".fisherman.yaml"),
        path.join(".fisherman.yml"),
        path.join(".fisherman.json"),
    ]
    .into_iter()
    .filter(|config| config.exists())
    .collect()
}

pub(crate) mod errors;

use crate::configuration::errors::ConfigurationError;
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
pub(crate) struct Configuration {
    pub hooks: Option<HashMap<GitHook, Vec<RuleRef>>>,
}

impl Configuration {
    pub(crate) fn load(path: &PathBuf) -> Result<Configuration, Box<dyn std::error::Error>> {
        let files = find_config_files(path.clone())?;

        let mut config = Figment::new();
        println!("{:?}", files);
        for file in files {
            let extension = match file.extension().and_then(OsStr::to_str) {
                Some(ext) => ext,
                None => return Err(Box::new(ConfigurationError::UnknownConfigFileExtension)),
            };

            config = match extension {
                "toml" => config.merge(Toml::file(file)),
                "yaml" => config.merge(Yaml::file(file)),
                "json" => config.merge(Json::file(file)),
                _ => return Err(Box::new(ConfigurationError::UnknownConfigFileExtension)),
            };
        }

        Ok(config.extract()?)
    }
}

fn find_config_files(path: PathBuf) -> Result<Vec<PathBuf>, Box<dyn std::error::Error>> {
    let locations = vec![home_dir().unwrap(), path.clone(), path.join(".git")];
    let mut config_files = vec![];
    for location_path in locations {
        let files = resolve_configs(location_path);
        if files.len() > 1 {
            return Err(Box::new(ConfigurationError::MultipleConfigFiles { files }));
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

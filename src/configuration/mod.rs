use crate::hooks::GitHook;
use dirs::home_dir;
use figment::providers::{Format, Json, Toml, Yaml};
use figment::Figment;
use serde::Deserialize;
use std::collections::HashMap;
use std::ffi::OsStr;
use std::fmt;
use std::path::PathBuf;

#[derive(Debug, Deserialize)]
#[serde(tag = "type")] // Use "tagged" enums to determine the variant
pub enum Item {
    // #[serde(rename = "A")] // Rename the variant to "A"
    VariantA { field_a: String },
    VariantB { field_b: i32 },
    VariantC { field_c: bool },
}

#[derive(Debug, Deserialize)]
pub(crate) struct Configuration {
    pub hooks: Option<HashMap<GitHook, Vec<Item>>>,
}

impl Default for Configuration {
    fn default() -> Self {
        Configuration { hooks: None }
    }
}

#[derive(Debug)]
pub(crate) enum ConfigurationError {
    MultipleConfigFiles { files: Vec<PathBuf> },
    UnknownConfigFileExtension,
}

impl fmt::Display for ConfigurationError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConfigurationError::MultipleConfigFiles { files } => {
                write!(f, "Multiple configuration files found: {:?}", files)
            }
            ConfigurationError::UnknownConfigFileExtension => {
                write!(f, "Unknown configuration file extension")
            }
        }
    }
}

impl std::error::Error for ConfigurationError {}

impl Configuration {
    pub(crate) fn load(path: &PathBuf) -> Result<Configuration, Box<dyn std::error::Error>> {
        let files = find_config_files(path.clone())?;

        let mut config = Figment::new();
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
    let locations = vec![path.clone(), path.join(".git"), home_dir().unwrap()];
    for path in locations {
        let configs = resolve_configs(path);
        if configs.len() > 1 {
            return Err(Box::new(ConfigurationError::MultipleConfigFiles {
                files: configs,
            }));
        } else if configs.len() == 1 {
            return Ok(configs);
        }
    }

    return Ok(vec![]);
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

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
        Self { hooks: None }
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
    let locations = vec![
        home_dir().unwrap(),
        path.clone(),
        path.join(".git")
    ];
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

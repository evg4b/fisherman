use crate::configuration::errors::ConfigurationError;
use anyhow::{bail, Result};
use dirs::home_dir;
use std::path::{Path, PathBuf};

pub fn find_config_files(path: &Path) -> Result<Vec<PathBuf>> {
    let locations = vec![path.join(".git"), path.to_path_buf(), home_dir().unwrap()];
    let mut config_files = vec![];
    for location_path in locations {
        let files = resolve_configs(location_path);
        match files.len() {
            0 => continue,
            1 => config_files.push(files[0].to_owned()),
            _ => bail!(ConfigurationError::MultipleConfigFiles { files }),
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

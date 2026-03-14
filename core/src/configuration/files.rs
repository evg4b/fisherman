use crate::configuration::errors::ConfigurationError;
use anyhow::{bail, Result};
use dirs::home_dir;
use std::path::{Path, PathBuf};

pub fn find_config_files(path: &Path) -> Result<Vec<PathBuf>> {
    // Now properly handles the case where home_dir() returns None by filtering it out of locations.
    let mut locations = vec![path.join(".git"), path.to_path_buf()];
    if let Some(home) = home_dir() {
        locations.push(home);
    }

    let mut config_files = vec![];
    for location_path in locations {
        let files = resolve_configs(location_path);
        match files.len() {
            0 => {}
            1 => config_files.push(files[0].clone()),
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

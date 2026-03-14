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

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use tempfile::TempDir;

    #[test]
    fn test_find_config_files_no_files() {
        let dir = TempDir::new().unwrap();
        let result = find_config_files(dir.path());
        assert!(result.is_ok());
    }

    #[test]
    fn test_find_config_files_single_toml() {
        let dir = TempDir::new().unwrap();
        let config_path = dir.path().join(".fisherman.toml");
        fs::write(&config_path, "[hooks]\n").unwrap();

        let result = find_config_files(dir.path()).unwrap();
        assert!(result.contains(&config_path));
    }

    #[test]
    fn test_find_config_files_single_yaml() {
        let dir = TempDir::new().unwrap();
        let config_path = dir.path().join(".fisherman.yaml");
        fs::write(&config_path, "hooks: {}\n").unwrap();

        let result = find_config_files(dir.path()).unwrap();
        assert!(result.contains(&config_path));
    }

    #[test]
    fn test_find_config_files_single_yml() {
        let dir = TempDir::new().unwrap();
        let config_path = dir.path().join(".fisherman.yml");
        fs::write(&config_path, "hooks: {}\n").unwrap();

        let result = find_config_files(dir.path()).unwrap();
        assert!(result.contains(&config_path));
    }

    #[test]
    fn test_find_config_files_single_json() {
        let dir = TempDir::new().unwrap();
        let config_path = dir.path().join(".fisherman.json");
        fs::write(&config_path, r#"{"hooks": {}}"#).unwrap();

        let result = find_config_files(dir.path()).unwrap();
        assert!(result.contains(&config_path));
    }

    #[test]
    fn test_find_config_files_multiple_in_path_errors() {
        let dir = TempDir::new().unwrap();
        fs::write(dir.path().join(".fisherman.toml"), "[hooks]\n").unwrap();
        fs::write(dir.path().join(".fisherman.yaml"), "hooks: {}\n").unwrap();

        let result = find_config_files(dir.path());
        assert!(result.is_err());
        assert!(result
            .unwrap_err()
            .to_string()
            .contains("Multiple config files found"));
    }

    #[test]
    fn test_find_config_files_multiple_in_git_dir_errors() {
        let dir = TempDir::new().unwrap();
        let git_dir = dir.path().join(".git");
        fs::create_dir(&git_dir).unwrap();
        fs::write(git_dir.join(".fisherman.toml"), "[hooks]\n").unwrap();
        fs::write(git_dir.join(".fisherman.yaml"), "hooks: {}\n").unwrap();

        let result = find_config_files(dir.path());
        assert!(result.is_err());
    }

    #[test]
    fn test_find_config_files_git_dir_config() {
        let dir = TempDir::new().unwrap();
        let git_dir = dir.path().join(".git");
        fs::create_dir(&git_dir).unwrap();
        let config_path = git_dir.join(".fisherman.toml");
        fs::write(&config_path, "[hooks]\n").unwrap();

        let result = find_config_files(dir.path()).unwrap();
        assert!(result.contains(&config_path));
    }
}

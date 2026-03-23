#![allow(dead_code)]

use std::fs;
use std::path::{Path, PathBuf};
use std::process::{Command, Output};
use tempdir::TempDir;

pub struct GitTestRepo {
    temp_dir: TempDir,
}

impl GitTestRepo {
    pub fn new() -> Self {
        let temp_dir = TempDir::new("fisherman_test").expect("Failed to create temp directory");
        let repo = Self { temp_dir };

        repo.git(&["init"]);
        repo.git(&["config", "user.name", "Test User"]);
        repo.git(&["config", "user.email", "test@example.com"]);
        repo.git(&["config", "commit.gpgsign", "false"]);

        repo
    }

    pub fn path(&self) -> &Path {
        self.temp_dir.path()
    }

    pub fn git(&self, args: &[&str]) -> Output {
        Command::new("git")
            .args(args)
            .current_dir(self.path())
            .output()
            .expect("Failed to execute git command")
    }

    pub fn create_file(&self, path: &str, content: &str) {
        let file_path = self.path().join(path);
        if let Some(parent) = file_path.parent() {
            fs::create_dir_all(parent).expect("Failed to create parent directories");
        }
        fs::write(file_path, content).expect("Failed to write file");
    }

    pub fn read_file(&self, path: &str) -> String {
        let file_path = self.path().join(path);
        fs::read_to_string(file_path).expect("Failed to read file")
    }

    pub fn file_exists(&self, path: &str) -> bool {
        self.path().join(path).exists()
    }

    pub fn create_config(&self, config: &str, format: ConfigFormat) {
        match format {
            ConfigFormat::Json => self.create_file(".fisherman.json", config),
            ConfigFormat::Yaml => self.create_file(".fisherman.yaml", config),
            ConfigFormat::Toml => self.create_file(".fisherman.toml", config),
        }
    }

    pub fn create_yaml_config(&self, config: &str) {
        self.create_file(".fisherman.yaml", config);
    }

    pub fn create_json_config(&self, config: &str) {
        self.create_file(".fisherman.json", config);
    }

    pub fn create_local_config(&self, config: &str) {
        self.create_file(".git/.fisherman.toml", config);
    }

    pub fn create_local_yaml_config(&self, config: &str) {
        self.create_file(".git/.fisherman.yaml", config);
    }

    pub fn create_local_json_config(&self, config: &str) {
        self.create_file(".git/.fisherman.json", config);
    }

    pub fn commit(&self, message: &str) -> Output {
        let add_output = self.git(&["add", "."]);
        assert!(
            add_output.status.success(),
            "Failed to stage files: {}",
            String::from_utf8_lossy(&add_output.stderr)
        );

        self.git(&["commit", "-m", message])
    }

    pub fn commit_allow_empty(&self, message: &str) -> Output {
        self.git(&["commit", "--allow-empty", "-m", message])
    }

    pub fn commit_with_hooks_allow_empty(&self, message: &str) -> Output {
        self.commit_allow_empty(message)
    }

    pub fn create_branch(&self, name: &str) {
        let output = self.git(&["checkout", "-b", name]);
        assert!(
            output.status.success(),
            "Failed to create branch: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn checkout(&self, name: &str) {
        let output = self.git(&["checkout", name]);
        assert!(
            output.status.success(),
            "Failed to checkout branch: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn checkout_new_branch(&self, name: &str) -> Output {
        self.git(&["checkout", "-b", name])
    }

    pub fn hook_exists(&self, hook_name: &str) -> bool {
        self.path().join(".git/hooks").join(hook_name).exists()
    }

    pub fn read_hook(&self, hook_name: &str) -> String {
        let hook_path = self.path().join(".git/hooks").join(hook_name);
        fs::read_to_string(hook_path).expect("Failed to read hook file")
    }

    pub fn commit_msg_file_path(&self) -> PathBuf {
        self.path().join(".git").join("COMMIT_EDITMSG")
    }

    pub fn write_commit_msg_file(&self, message: &str) {
        let path = self.commit_msg_file_path();
        fs::write(&path, message).expect("Failed to write commit message file");
    }

    pub fn git_history(&self, commits: &[(&str, &[(&str, &str)])]) {
        for (message, files) in commits {
            for (path, content) in *files {
                self.create_file(path, content);
            }
            let output = self.commit(message);
            assert!(
                output.status.success(),
                "Failed to create commit '{}': {}",
                message,
                String::from_utf8_lossy(&output.stderr)
            );
        }
    }
}

impl Default for GitTestRepo {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ConfigFormat {
    Toml,
    Yaml,
    Json,
}

pub struct ConfigBuilder<'a> {
    repo: &'a mut GitTestRepo,
    configs: Vec<(ConfigFormat, String, bool)>,
}

impl<'a> ConfigBuilder<'a> {
    pub fn new(repo: &'a mut GitTestRepo) -> Self {
        Self {
            repo,
            configs: Vec::new(),
        }
    }

    pub fn repository(mut self, content: &str) -> Self {
        self.configs
            .push((ConfigFormat::Toml, content.to_string(), false));
        self
    }

    pub fn repository_with_format(mut self, format: ConfigFormat, content: &str) -> Self {
        self.configs.push((format, content.to_string(), false));
        self
    }

    pub fn local(mut self, content: &str) -> Self {
        self.configs
            .push((ConfigFormat::Toml, content.to_string(), true));
        self
    }

    pub fn local_with_format(mut self, format: ConfigFormat, content: &str) -> Self {
        self.configs.push((format, content.to_string(), true));
        self
    }

    pub fn repository_config(self, config: &fisherman_core::Configuration) -> Self {
        self.repository_with_format(
            ConfigFormat::Toml,
            &crate::common::configuration::serialize_configuration(config, ConfigFormat::Toml),
        )
    }

    pub fn repository_config_with_format(
        self,
        config: &fisherman_core::Configuration,
        format: ConfigFormat,
    ) -> Self {
        self.repository_with_format(
            format,
            &crate::common::configuration::serialize_configuration(config, format),
        )
    }

    pub fn local_config(self, config: &fisherman_core::Configuration) -> Self {
        self.local_with_format(
            ConfigFormat::Toml,
            &crate::common::configuration::serialize_configuration(config, ConfigFormat::Toml),
        )
    }

    pub fn local_config_with_format(
        self,
        config: &fisherman_core::Configuration,
        format: ConfigFormat,
    ) -> Self {
        self.local_with_format(
            format,
            &crate::common::configuration::serialize_configuration(config, format),
        )
    }

    pub fn build(self) {
        for (format, content, is_local) in self.configs {
            if is_local {
                match format {
                    ConfigFormat::Toml => self.repo.create_local_config(&content),
                    ConfigFormat::Yaml => self.repo.create_local_yaml_config(&content),
                    ConfigFormat::Json => self.repo.create_local_json_config(&content),
                }
            } else {
                self.repo.create_config(&content, format);
            }
        }
    }
}

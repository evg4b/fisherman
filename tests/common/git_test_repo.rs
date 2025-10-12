#![allow(dead_code)]

use std::env;
use std::fs;
use std::path::{Path, PathBuf};
use std::process::{Command, Output};
use tempdir::TempDir;

pub struct GitTestRepo {
    temp_dir: TempDir,
    global_config_dir: Option<TempDir>,
}

impl GitTestRepo {
    pub fn new() -> Self {
        let temp_dir = TempDir::new("fisherman_test").expect("Failed to create temp directory");
        let repo = Self {
            temp_dir,
            global_config_dir: None,
        };

        repo.git(&["init"]);
        repo.git(&["config", "user.name", "Test User"]);
        repo.git(&["config", "user.email", "test@example.com"]);
        repo.git(&["config", "commit.gpgsign", "false"]);

        repo
    }

    /// Enable global config support by creating a temporary home directory
    pub fn with_global_config(&mut self) -> &mut Self {
        if self.global_config_dir.is_none() {
            let global_dir = TempDir::new("fisherman_global").expect("Failed to create global config dir");
            self.global_config_dir = Some(global_dir);
        }
        self
    }

    /// Get the global config directory path (creates it if needed)
    pub fn global_config_path(&mut self) -> PathBuf {
        self.with_global_config();
        self.global_config_dir.as_ref().unwrap().path().to_path_buf()
    }

    /// Set HOME environment variable to point to our test global config dir
    pub fn use_global_config(&mut self) {
        self.with_global_config();
        let global_path = self.global_config_dir.as_ref().unwrap().path();
        // SAFETY: This is only used in tests to set a temporary HOME directory.
        // Tests run sequentially in the same process, and we're setting it to a valid path.
        unsafe {
            env::set_var("HOME", global_path);
        }
    }

    /// Create a global config file
    pub fn create_global_config(&mut self, config: &str) {
        let global_path = self.global_config_path();
        let config_path = global_path.join(".fisherman.toml");
        fs::write(config_path, config).expect("Failed to write global config");
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

    pub fn create_config(&self, config: &str) {
        self.create_file(".fisherman.toml", config);
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

    pub fn current_branch(&self) -> String {
        let output = self.git(&["branch", "--show-current"]);
        String::from_utf8_lossy(&output.stdout).trim().to_string()
    }

    pub fn hook_exists(&self, hook_name: &str) -> bool {
        self.path().join(".git/hooks").join(hook_name).exists()
    }

    pub fn read_hook(&self, hook_name: &str) -> String {
        let hook_path = self.path().join(".git/hooks").join(hook_name);
        fs::read_to_string(hook_path).expect("Failed to read hook file")
    }

    pub fn write_commit_msg_file(&self, message: &str) {
        let msg_path = self.path().join(".git/COMMIT_EDITMSG");
        fs::write(msg_path, message).expect("Failed to write commit message file");
    }

    pub fn commit_msg_file_path(&self) -> PathBuf {
        self.path().join(".git/COMMIT_EDITMSG")
    }

    /// Creates a Git history with multiple commits and files
    ///
    /// # Example
    /// ```
    /// repo.git_history(&[
    ///     ("Initial commit", &[
    ///         ("README.md", "# Project"),
    ///         ("src/main.rs", "fn main() {}"),
    ///     ]),
    ///     ("Add tests", &[
    ///         ("tests/test.rs", "#[test] fn test() {}"),
    ///     ]),
    /// ]);
    /// ```
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

/// Represents configuration scope levels
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ConfigScope {
    Global,
    Repository,
    Local,
}

/// Represents configuration file format
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ConfigFormat {
    Toml,
    Yaml,
    Json,
}

/// Builder for creating multiple scoped configurations in tests
pub struct ConfigBuilder<'a> {
    repo: &'a mut GitTestRepo,
    configs: Vec<(ConfigScope, ConfigFormat, String)>,
}

impl<'a> ConfigBuilder<'a> {
    pub fn new(repo: &'a mut GitTestRepo) -> Self {
        Self {
            repo,
            configs: Vec::new(),
        }
    }

    /// Add a global config
    pub fn global(mut self, content: &str) -> Self {
        self.configs.push((ConfigScope::Global, ConfigFormat::Toml, content.to_string()));
        self
    }

    /// Add a repository config (default scope)
    pub fn repository(mut self, content: &str) -> Self {
        self.configs.push((ConfigScope::Repository, ConfigFormat::Toml, content.to_string()));
        self
    }

    /// Add a repository config with specific format
    pub fn repository_with_format(mut self, format: ConfigFormat, content: &str) -> Self {
        self.configs.push((ConfigScope::Repository, format, content.to_string()));
        self
    }

    /// Add a local config (.git/.fisherman.toml)
    pub fn local(mut self, content: &str) -> Self {
        self.configs.push((ConfigScope::Local, ConfigFormat::Toml, content.to_string()));
        self
    }

    /// Build all configurations
    pub fn build(self) {
        for (scope, format, content) in self.configs {
            match (scope, format) {
                (ConfigScope::Global, ConfigFormat::Toml) => {
                    self.repo.create_global_config(&content);
                }
                (ConfigScope::Repository, ConfigFormat::Toml) => {
                    self.repo.create_config(&content);
                }
                (ConfigScope::Repository, ConfigFormat::Yaml) => {
                    self.repo.create_yaml_config(&content);
                }
                (ConfigScope::Repository, ConfigFormat::Json) => {
                    self.repo.create_json_config(&content);
                }
                (ConfigScope::Local, ConfigFormat::Toml) => {
                    self.repo.create_local_config(&content);
                }
                _ => panic!("Unsupported config scope/format combination: {:?}/{:?}", scope, format),
            }
        }
    }
}

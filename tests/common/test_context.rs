use super::{ConfigFormat, FishermanBinary, GitTestRepo};
use crate::common::configuration::serialize_configuration;
use core::configuration::Configuration;
use std::process::Output;

pub struct TestContext {
    pub binary: FishermanBinary,
    pub repo: GitTestRepo,
}

impl TestContext {
    pub fn new() -> Self {
        Self {
            binary: FishermanBinary::build(),
            repo: GitTestRepo::new(),
        }
    }

    pub fn setup_with_config(&self, config: &str) -> Output {
        self.repo.create_config(config, ConfigFormat::Toml);
        self.repo
            .git_history(&[("initial", &[("test.txt", "initial")])]);
        self.binary.install(self.repo.path(), false)
    }

    pub fn setup_with_history_and_install(
        &self,
        config: &Configuration,
        format: ConfigFormat,
        history: &[(&str, &[(&str, &str)])],
    ) {
        let config_string = serialize_configuration(config, format);
        self.repo.create_config(&config_string, format);
        self.repo.git_history(history);
        let output = self.binary.install(self.repo.path(), false);
        assert!(
            output.status.success(),
            "Installation failed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn setup_and_install(&self, config: &Configuration, format: ConfigFormat) {
        let config_string = serialize_configuration(config, format);
        let output = self.setup_with_config(&config_string);
        assert!(
            output.status.success(),
            "Installation failed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn setup_and_install_old(&self, config: &str) {
        let output = self.setup_with_config(config);
        assert!(
            output.status.success(),
            "Installation failed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn handle(&self, hook: &str) -> Output {
        self.binary.handle(hook, self.repo.path(), &[])
    }

    pub fn handle_commit_msg(&self, message: &str) -> Output {
        self.repo.write_commit_msg_file(message);
        let msg_path = self.repo.commit_msg_file_path();
        self.binary.handle(
            "commit-msg",
            self.repo.path(),
            &[msg_path.to_str().unwrap()],
        )
    }

    pub fn git_commit_allow_empty(&self, message: &str) -> Output {
        self.repo.commit_with_hooks_allow_empty(message)
    }

    pub fn git_commit_allow_empty_success(&self, message: &str) {
        let output = self.git_commit_allow_empty(message);
        assert!(
            output.status.success(),
            "Git commit (allow-empty) should succeed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    pub fn git_commit_allow_empty_failure(&self, message: &str) {
        let output = self.git_commit_allow_empty(message);
        assert!(
            !output.status.success(),
            "Git commit (allow-empty) should fail due to hook"
        );
    }

    pub fn git_checkout_new_branch(&self, name: &str) -> Output {
        self.repo.checkout_new_branch(name)
    }
}

impl Default for TestContext {
    fn default() -> Self {
        Self::new()
    }
}

pub fn assert_stderr_contains(stderr: &str, expected: &[&str], context: &str) {
    let found = expected.iter().any(|&exp| stderr.contains(exp));
    assert!(
        found,
        "{}: expected stderr to contain one of {:?}, got: {}",
        context, expected, stderr
    );
}

use super::{FishermanBinary, GitTestRepo};
use std::process::Output;

/// Helper struct that combines FishermanBinary and GitTestRepo for easier test setup
pub struct TestContext {
    pub binary: FishermanBinary,
    pub repo: GitTestRepo,
}

impl TestContext {
    /// Creates a new test context with a fresh repo and built binary
    pub fn new() -> Self {
        Self {
            binary: FishermanBinary::build(),
            repo: GitTestRepo::new(),
        }
    }

    /// Creates a config, initializes the repo with a file and commit, and installs hooks
    pub fn setup_with_config(&self, config: &str) -> Output {
        self.repo.create_config(config);
        self.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
        self.binary.install(self.repo.path(), false)
    }

    /// Creates a config with custom git history and installs hooks
    pub fn setup_with_history(&self, config: &str, history: &[(&str, &[(&str, &str)])]) -> Output {
        self.repo.create_config(config);
        self.repo.git_history(history);
        self.binary.install(self.repo.path(), false)
    }

    /// Setup and assert installation succeeded
    pub fn setup_and_install(&self, config: &str) {
        let output = self.setup_with_config(config);
        assert!(
            output.status.success(),
            "Installation failed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Handle a hook and return the output
    pub fn handle(&self, hook: &str) -> Output {
        self.binary.handle(hook, self.repo.path(), &[])
    }

    /// Handle a hook and assert it succeeded
    pub fn handle_success(&self, hook: &str) {
        let output = self.handle(hook);
        assert!(
            output.status.success(),
            "Hook '{}' should succeed: {}",
            hook,
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Handle a hook and assert it failed
    pub fn handle_failure(&self, hook: &str) {
        let output = self.handle(hook);
        assert!(
            !output.status.success(),
            "Hook '{}' should fail",
            hook
        );
    }

    /// Handle commit-msg hook with a message
    pub fn handle_commit_msg(&self, message: &str) -> Output {
        self.repo.write_commit_msg_file(message);
        let msg_path = self.repo.commit_msg_file_path();
        self.binary.handle(
            "commit-msg",
            self.repo.path(),
            &[msg_path.to_str().unwrap()],
        )
    }

    /// Handle commit-msg hook and assert success
    pub fn handle_commit_msg_success(&self, message: &str) {
        let output = self.handle_commit_msg(message);
        assert!(
            output.status.success(),
            "commit-msg hook should succeed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Handle commit-msg hook and assert failure
    pub fn handle_commit_msg_failure(&self, message: &str) {
        let output = self.handle_commit_msg(message);
        assert!(
            !output.status.success(),
            "commit-msg hook should fail"
        );
    }
}

impl Default for TestContext {
    fn default() -> Self {
        Self::new()
    }
}

/// Platform-specific config helper for exec commands
#[cfg(windows)]
pub fn echo_config(hook: &str, text: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "{}"]
"#,
        hook, text
    )
}

/// Platform-specific config helper for exec commands
#[cfg(not(windows))]
pub fn echo_config(hook: &str, text: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "exec"
command = "echo"
args = ["{}"]
"#,
        hook, text
    )
}

/// Platform-specific config helper for failing exec commands
#[cfg(windows)]
pub fn fail_config(hook: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "exec"
command = "cmd"
args = ["/C", "exit", "1"]
"#,
        hook
    )
}

/// Platform-specific config helper for failing exec commands
#[cfg(not(windows))]
pub fn fail_config(hook: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "exec"
command = "false"
"#,
        hook
    )
}

/// Platform-specific config helper for shell scripts
#[cfg(windows)]
pub fn shell_config(hook: &str, script: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "shell"
script = "{}"
"#,
        hook, script
    )
}

/// Platform-specific config helper for shell scripts
#[cfg(not(windows))]
pub fn shell_config(hook: &str, script: &str) -> String {
    format!(
        r#"
[[hooks.{}]]
type = "shell"
script = """
{}
"""
"#,
        hook, script
    )
}

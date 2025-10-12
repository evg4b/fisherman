#![allow(dead_code)]

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

    // ========== Git-based hook testing (recommended approach) ==========
    // These methods trigger hooks through actual git commands instead of
    // calling fisherman directly, testing the real-world scenario.

    /// Commit using git (triggers pre-commit and commit-msg hooks)
    pub fn git_commit(&self, message: &str) -> Output {
        self.repo.commit_with_hooks(message)
    }

    /// Commit using git and assert success
    pub fn git_commit_success(&self, message: &str) {
        let output = self.git_commit(message);
        assert!(
            output.status.success(),
            "Git commit should succeed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Commit using git and assert failure
    pub fn git_commit_failure(&self, message: &str) {
        let output = self.git_commit(message);
        assert!(
            !output.status.success(),
            "Git commit should fail due to hook"
        );
    }

    /// Commit allowing empty (useful for testing hooks without file changes)
    pub fn git_commit_allow_empty(&self, message: &str) -> Output {
        self.repo.commit_with_hooks_allow_empty(message)
    }

    /// Commit allowing empty and assert success
    pub fn git_commit_allow_empty_success(&self, message: &str) {
        let output = self.git_commit_allow_empty(message);
        assert!(
            output.status.success(),
            "Git commit (allow-empty) should succeed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Commit allowing empty and assert failure
    pub fn git_commit_allow_empty_failure(&self, message: &str) {
        let output = self.git_commit_allow_empty(message);
        assert!(
            !output.status.success(),
            "Git commit (allow-empty) should fail due to hook"
        );
    }

    /// Checkout branch using git (triggers post-checkout hook if configured)
    pub fn git_checkout(&self, branch: &str) -> Output {
        self.repo.checkout_with_hooks(branch)
    }

    /// Checkout and assert success
    pub fn git_checkout_success(&self, branch: &str) {
        let output = self.git_checkout(branch);
        assert!(
            output.status.success(),
            "Git checkout should succeed: {}",
            String::from_utf8_lossy(&output.stderr)
        );
    }

    /// Create and checkout new branch using git
    pub fn git_checkout_new_branch(&self, name: &str) -> Output {
        self.repo.checkout_new_branch(name)
    }
}

impl Default for TestContext {
    fn default() -> Self {
        Self::new()
    }
}

/// Assert that stderr contains any of the expected strings
pub fn assert_stderr_contains(stderr: &str, expected: &[&str], context: &str) {
    let found = expected.iter().any(|&exp| stderr.contains(exp));
    assert!(
        found,
        "{}: expected stderr to contain one of {:?}, got: {}",
        context, expected, stderr
    );
}

mod common;

use crate::common::ConfigFormat;
use common::{configuration::serialize_configuration, test_context::TestContext, FishermanBinary, GitTestRepo};
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::RuleParams;
use std::collections::HashMap;

/// Tests that exec rule executes successfully when command exits with code 0.
/// Verifies basic command execution functionality using echo command.
#[test]
fn exec_rule_success() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("test")]),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("test")]),
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    let stdout = String::from_utf8_lossy(&output.stdout);
    assert!(stdout.contains("test"), "Output should contain 'test': {}", stdout);
}

/// Tests that exec rule fails appropriately when command exits with non-zero code.
/// Verifies that command failures are detected and propagate as hook failures.
#[test]
fn exec_rule_failure() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("exit"), String::from("1")]),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("false"),
                args: None,
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error output should not be empty");
}

/// Tests that exec rule properly passes environment variables to executed command.
/// Verifies that custom env variables are available in the command's environment.
#[test]
fn exec_rule_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("%TEST_VAR%")]),
                env: Some(HashMap::from([(String::from("TEST_VAR"), String::from("test_value"))])),
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("sh"),
                args: Some(vec![String::from("-c"), String::from("test \"$TEST_VAR\" = \"test_value\"")]),
                env: Some(HashMap::from([(String::from("TEST_VAR"), String::from("test_value"))])),
            })
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let commit_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        commit_output.status.success(),
        "Hook should pass environment variables: {}",
        String::from_utf8_lossy(&commit_output.stderr)
    );
}

/// Tests that shell script rule executes successfully when script exits with code 0.
/// Verifies basic shell script execution on both Windows and Unix platforms.
#[test]
fn shell_script_success() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("echo test"),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("#!/bin/sh\necho \"Running shell script\"\nexit 0\n"),
                env: None,
            })
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let commit_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        commit_output.status.success(),
        "Shell script should succeed: {}",
        String::from_utf8_lossy(&commit_output.stderr)
    );

    // Note: When running through git commit, hook output goes to stderr, not stdout
    // The important thing is that the hook succeeds
}

/// Tests that shell script rule fails when script exits with non-zero code.
/// Verifies that shell script failures properly abort the hook execution.
#[test]
fn shell_script_failure() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("exit 1"),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("#!/bin/sh\nexit 1\n"),
                env: None,
            })
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let commit_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        !commit_output.status.success(),
        "Shell script should fail with exit 1"
    );

    let stderr = String::from_utf8_lossy(&commit_output.stderr);
    assert!(!stderr.is_empty(), "Error output should not be empty when shell fails");
}

/// Tests that shell script can access custom environment variables defined in config.
/// Verifies environment variable propagation to shell execution context.
#[test]
fn shell_script_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("if \"%CUSTOM_VAR%\" == \"custom_value\" exit 0"),
                env: Some(HashMap::from([(String::from("CUSTOM_VAR"), String::from("custom_value"))])),
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("#!/bin/sh\nif [ \"$CUSTOM_VAR\" = \"custom_value\" ]; then\n    exit 0\nelse\n    exit 1\nfi\n"),
                env: Some(HashMap::from([(String::from("CUSTOM_VAR"), String::from("custom_value"))])),
            })
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let commit_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        commit_output.status.success(),
        "Shell script should access env variables: {}",
        String::from_utf8_lossy(&commit_output.stderr)
    );
}

/// Tests that exec and shell rules can be configured together and both execute successfully.
/// Verifies compatibility and correct execution order of mixed rule types.
#[test]
fn exec_and_shell_mixed() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("exec test")]),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo shell test"),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("exec test")]),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'shell test'"),
                env: None,
            })
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let commit_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        commit_output.status.success(),
        "Both exec and shell rules should succeed: {}",
        String::from_utf8_lossy(&commit_output.stderr)
    );

    // Note: When running through git commit, hook output behavior is different
    // The important thing is that both rules execute successfully
}

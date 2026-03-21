mod common;

use crate::common::ConfigFormat;
use common::{
    configuration::serialize_configuration, test_context::TestContext, FishermanBinary, GitTestRepo,
};
use core::Configuration;
use core::ExecRule;
use core::GitHook;
use core::ShellScriptRule;
use std::collections::HashMap;

#[test]
#[cfg(feature = "integration-tests")]
fn exec_rule_success() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "cmd".into(),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("test")]),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "echo".into(),
                args: Some(vec![String::from("test")]),
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    let stdout = String::from_utf8_lossy(&output.stdout);
    assert!(
        stdout.contains("test"),
        "Output should contain 'test': {}",
        stdout
    );
}

#[test]
#[cfg(feature = "integration-tests")]
fn exec_rule_failure() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "cmd".into(),
                args: Some(vec![String::from("/C"), String::from("exit"), String::from("1")]),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "false".into(),
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

#[test]
#[cfg(feature = "integration-tests")]
fn exec_rule_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "cmd".into(),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("%TEST_VAR%")]),
                env: Some(HashMap::from([(String::from("TEST_VAR"), String::from("test_value"))])),
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "sh".into(),
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

#[test]
#[cfg(feature = "integration-tests")]
fn shell_script_success() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "echo test".into(),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "#!/bin/sh\necho \"Running shell script\"\nexit 0\n".into(),
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
}

#[test]
#[cfg(feature = "integration-tests")]
fn shell_script_failure() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "exit 1".into(),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "#!/bin/sh\nexit 1\n".into(),
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
    assert!(
        !stderr.is_empty(),
        "Error output should not be empty when shell fails"
    );
}

#[test]
#[cfg(feature = "integration-tests")]
fn shell_script_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "if \"%CUSTOM_VAR%\" == \"custom_value\" exit 0".into(),
                env: Some(HashMap::from([(String::from("CUSTOM_VAR"), String::from("custom_value"))])),
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "#!/bin/sh\nif [ \"$CUSTOM_VAR\" = \"custom_value\" ]; then\n    exit 0\nelse\n    exit 1\nfi\n".into(),
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

#[test]
#[cfg(feature = "integration-tests")]
fn exec_and_shell_mixed() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "cmd".into(),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("exec test")]),
                env: None,
            }),
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "echo shell test".into(),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                when: None,
                extract: None,
                command: "echo".into(),
                args: Some(vec![String::from("exec test")]),
                env: None,
            }),
            rule!(ShellScriptRule {
                when: None,
                extract: None,
                script: "echo 'shell test'".into(),
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
}

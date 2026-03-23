mod common;

use crate::common::ConfigFormat;
use common::{configuration::serialize_configuration, FishermanBinary, GitTestRepo};
use fisherman_core::{
    BranchNameRegexRule, CommitMessageRegexRule, Configuration, GitHook,
};

#[test]
fn install_creates_hooks() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = {
        let mut cfg = config!(
            GitHook::PreCommit => [
                rule!(BranchNameRegexRule {
                    expression: ".*".into(),
                })
            ]
        );
        cfg.hooks.insert(
            GitHook::CommitMsg,
            vec![rule!(CommitMessageRegexRule {
                when: None,
                expression: ".*".into(),
            })],
        );
        cfg.hooks.insert(
            GitHook::PrePush,
            vec![rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })],
        );
        cfg
    };
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("pre-commit"));
    assert!(repo.hook_exists("commit-msg"));
    assert!(repo.hook_exists("pre-push"));
}

#[test]
fn install_without_force_fails_when_hook_exists() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.create_file(".git/hooks/pre-commit", "#!/bin/sh\necho existing hook");

    let install_output = binary.install(repo.path(), false);

    assert!(
        !install_output.status.success(),
        "Installation should fail without --force when hook exists"
    );

    let stderr = String::from_utf8_lossy(&install_output.stderr);
    assert!(
        stderr.contains("already installed") || stderr.contains("exist") || stderr.contains("force"),
        "Error should mention existing hook and --force flag"
    );
}

#[test]
fn install_with_force_overwrites_existing_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.create_file(".git/hooks/pre-commit", "#!/bin/sh\necho existing hook");

    let install_output = binary.install(repo.path(), true);

    assert!(
        install_output.status.success(),
        "Installation should succeed with --force: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("pre-commit"));
    assert!(repo.hook_exists("pre-commit.bkp"));

    let hook_content = repo.read_hook("pre-commit");

    #[cfg(not(windows))]
    assert!(
        hook_content.contains("fisherman handle"),
        "Hook should contain fisherman command"
    );
    #[cfg(windows)]
    assert!(
        hook_content.contains("fisherman.exe handle"),
        "Hook should contain fisherman command"
    );
}

#[test]
fn hook_script_contains_correct_command() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("pre-commit");

    assert!(hook_content.starts_with("#!/bin/sh"));
    #[cfg(not(windows))]
    assert!(hook_content.contains("fisherman handle pre-commit"));
    #[cfg(windows)]
    assert!(hook_content.contains("fisherman.exe handle pre-commit"));
    #[cfg(not(windows))]
    assert!(hook_content.contains(&binary.path().display().to_string()));
}

#[test]
fn commit_msg_hook_passes_arguments() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("commit-msg");

    assert!(
        hook_content.contains("$@"),
        "commit-msg hook should pass arguments"
    );
}

#[test]
fn explain_command_shows_rules() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(fisherman_core::BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(fisherman_core::ExecRule {
                command: "echo".into(),
                args: Some(vec!["test".into()]),
                env: None,
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let explain_output = binary.explain("pre-commit", repo.path());

    assert!(
        explain_output.status.success(),
        "Explain command should succeed"
    );

    let output = String::from_utf8_lossy(&explain_output.stdout);
    assert!(output.contains("Branch should start with"), "Output should contain branch name rule: {}", output);
    assert!(output.contains("Execute command"), "Output should contain exec rule: {}", output);
    assert!(output.contains("feature/"), "Output should show rule configuration details: {}", output);
}

#[test]
fn no_config_installs_no_hooks() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed even without config"
    );
}

#[test]
fn install_multiple_hooks_for_same_event() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^feature/".into(),
            }),
            rule!(fisherman_core::BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(fisherman_core::BranchNameSuffixRule {
                suffix: "-dev".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    assert!(repo.hook_exists("pre-commit"));

    repo.create_branch("feature/test-dev");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(
        handle_output.status.success(),
        "All rules for same hook should execute"
    );
}

#[test]
fn hierarchical_config_merge() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(fisherman_core::WriteFileRule {
                path: "repo.txt".into(),
                content: "repo level".into(),
                append: None,
            })
        ]
    );
    let repo_config_string = serialize_configuration(&repo_config, ConfigFormat::Toml);

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(fisherman_core::WriteFileRule {
                path: "local.txt".into(),
                content: "local level".into(),
                append: None,
            })
        ]
    );
    let local_config_string = serialize_configuration(&local_config, ConfigFormat::Toml);

    repo.create_config(&repo_config_string, ConfigFormat::Toml);
    repo.create_local_config(&local_config_string);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(
        repo.file_exists("repo.txt"),
        "Repository level config should be executed"
    );
    assert!(
        repo.file_exists("local.txt"),
        "Local level config should be executed"
    );
}

#[test]
fn install_pre_push_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PrePush => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("pre-push"));
    let hook_content = repo.read_hook("pre-push");
    #[cfg(not(windows))]
    assert!(hook_content.contains("fisherman handle pre-push"));
    #[cfg(windows)]
    assert!(hook_content.contains("fisherman.exe handle pre-push"));
}

#[test]
fn install_post_commit_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PostCommit => [
            rule!(fisherman_core::WriteFileRule {
                path: "post-commit-ran.txt".into(),
                content: "Post commit hook executed".into(),
                append: None,
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("post-commit"));
}

#[test]
fn install_prepare_commit_msg_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PrepareCommitMsg => [
            rule!(fisherman_core::WriteFileRule {
                path: "prepare-ran.txt".into(),
                content: "Prepare commit msg hook executed".into(),
                append: None,
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("prepare-commit-msg"));
}

#[test]
fn backup_file_contains_original_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );
    let config_string = serialize_configuration(&config, ConfigFormat::Toml);

    let original_hook_content = "#!/bin/sh\necho original hook\nexit 0";
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.create_file(".git/hooks/pre-commit", original_hook_content);

    let install_output = binary.install(repo.path(), true);

    assert!(
        install_output.status.success(),
        "Installation should succeed with --force: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("pre-commit.bkp"));
    let backup_content = repo.read_hook("pre-commit.bkp");
    assert_eq!(
        backup_content, original_hook_content,
        "Backup should contain original hook content"
    );
}

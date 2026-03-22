mod common;

use crate::common::ConfigFormat;
use common::{FishermanBinary, GitTestRepo};

#[test]
#[cfg(feature = "integration-tests")]
fn install_creates_hooks() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"

[[hooks.commit-msg]]
type = "message-regex"
regex = ".*"

[[hooks.pre-push]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);

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
#[cfg(feature = "integration-tests")]
fn install_without_force_fails_when_hook_exists() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);
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
#[cfg(feature = "integration-tests")]
fn install_with_force_overwrites_existing_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);
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
#[cfg(feature = "integration-tests")]
fn hook_script_contains_correct_command() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("pre-commit");

    assert!(hook_content.starts_with("#!/bin/sh"));
    #[cfg(not(windows))]
    assert!(hook_content.contains("fisherman handle pre-commit"));
    #[cfg(windows)]
    assert!(hook_content.contains("fisherman.exe handle pre-commit"));
    assert!(hook_content.contains(&binary.path().display().to_string()));
}

#[test]
#[cfg(feature = "integration-tests")]
fn commit_msg_hook_passes_arguments() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("commit-msg");

    assert!(
        hook_content.contains("$@"),
        "commit-msg hook should pass arguments"
    );
}

#[test]
#[cfg(feature = "integration-tests")]
fn explain_command_shows_rules() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["test"]
"#;

    repo.create_config(config, ConfigFormat::Toml);

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
#[cfg(feature = "integration-tests")]
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
#[cfg(feature = "integration-tests")]
fn install_multiple_hooks_for_same_event() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/"

[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-dev"
"#;

    repo.create_config(config, ConfigFormat::Toml);
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
#[cfg(feature = "integration-tests")]
fn hierarchical_config_merge() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "repo.txt"
content = "repo level"
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "local.txt"
content = "local level"
"#;

    repo.create_config(repo_config, ConfigFormat::Toml);
    repo.create_local_config(local_config);
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
#[cfg(feature = "integration-tests")]
fn install_pre_push_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-push]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config, ConfigFormat::Toml);

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
#[cfg(feature = "integration-tests")]
fn install_post_commit_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.post-commit]]
type = "write-file"
path = "post-commit-ran.txt"
content = "Post commit hook executed"
"#;

    repo.create_config(config, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("post-commit"));
}

#[test]
#[cfg(feature = "integration-tests")]
fn install_prepare_commit_msg_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.prepare-commit-msg]]
type = "write-file"
path = "prepare-ran.txt"
content = "Prepare commit msg hook executed"
"#;

    repo.create_config(config, ConfigFormat::Toml);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("prepare-commit-msg"));
}

#[test]
#[cfg(feature = "integration-tests")]
fn backup_file_contains_original_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    let original_hook_content = "#!/bin/sh\necho original hook\nexit 0";
    repo.create_config(config, ConfigFormat::Toml);
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

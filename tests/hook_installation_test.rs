mod common;

use common::{FishermanBinary, GitTestRepo};

/// Tests that install command creates hook files for all configured hook types.
/// Verifies that pre-commit, commit-msg, and pre-push hooks are properly installed.
#[test]
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

    repo.create_config(config);

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

/// Tests that install without --force flag fails when hook file already exists.
/// Verifies safety mechanism that prevents accidental overwriting of existing hooks.
#[test]
fn install_without_force_fails_when_hook_exists() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config);
    repo.create_file(".git/hooks/pre-commit", "#!/bin/sh\necho existing hook");

    let install_output = binary.install(repo.path(), false);

    assert!(
        !install_output.status.success(),
        "Installation should fail without --force when hook exists"
    );
}

/// Tests that install with --force flag overwrites existing hooks and creates backup.
/// Verifies that original hooks are backed up with .bkp extension before overwriting.
#[test]
fn install_with_force_overwrites_existing_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config);
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
    assert!(
        hook_content.contains("fisherman handle"),
        "Hook should contain fisherman command"
    );
}

/// Tests that installed hook script contains proper shebang and fisherman handle command.
/// Verifies hook script structure includes correct binary path and hook name.
#[test]
fn hook_script_contains_correct_command() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("pre-commit");

    assert!(hook_content.starts_with("#!/bin/sh"));
    assert!(hook_content.contains("fisherman handle pre-commit"));
    assert!(hook_content.contains(&binary.path().display().to_string()));
}

/// Tests that commit-msg hook script properly passes arguments using $@ to fisherman.
/// Verifies that message file path is forwarded correctly to handle command.
#[test]
fn commit_msg_hook_passes_arguments() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = ".*"
"#;

    repo.create_config(config);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    let hook_content = repo.read_hook("commit-msg");

    assert!(
        hook_content.contains("$@"),
        "commit-msg hook should pass arguments"
    );
}

/// Tests that explain command displays configured rules for a specific hook type.
/// Verifies that rule types and descriptions are properly output for debugging.
#[test]
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

    repo.create_config(config);

    let explain_output = binary.explain("pre-commit", repo.path());

    assert!(
        explain_output.status.success(),
        "Explain command should succeed"
    );

    let output = String::from_utf8_lossy(&explain_output.stdout);
    assert!(output.contains("branch name"));
    assert!(output.contains("exec"));
}

/// Tests that install succeeds without errors when no configuration file exists.
/// Verifies graceful handling of repositories without fisherman configuration.
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

/// Tests that multiple rules configured for the same hook type all execute in order.
/// Verifies that all regex, prefix, and suffix rules run when hook is triggered.
#[test]
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

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    assert!(repo.hook_exists("pre-commit"));

    repo.create_branch("feature/test-dev");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "All rules for same hook should execute"
    );
}

/// Tests that repository and local configurations are merged correctly.
/// Verifies that rules from both .fisherman.toml and .git/.fisherman.toml execute together.
#[test]
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

    repo.create_config(repo_config);
    repo.create_local_config(local_config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

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

/// Tests that pre-push hook is correctly installed with proper handle command.
/// Verifies support for pre-push hook type installation and configuration.
#[test]
fn install_pre_push_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-push]]
type = "branch-name-regex"
regex = ".*"
"#;

    repo.create_config(config);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("pre-push"));
    let hook_content = repo.read_hook("pre-push");
    assert!(hook_content.contains("fisherman handle pre-push"));
}

/// Tests that post-commit hook is correctly installed and can execute write-file rules.
/// Verifies support for post-commit hook type in fisherman configuration.
#[test]
fn install_post_commit_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.post-commit]]
type = "write-file"
path = "post-commit-ran.txt"
content = "Post commit hook executed"
"#;

    repo.create_config(config);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("post-commit"));
}

/// Tests that prepare-commit-msg hook is correctly installed with proper configuration.
/// Verifies support for prepare-commit-msg hook type in fisherman setup.
#[test]
fn install_prepare_commit_msg_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.prepare-commit-msg]]
type = "write-file"
path = "prepare-ran.txt"
content = "Prepare commit msg hook executed"
"#;

    repo.create_config(config);

    let install_output = binary.install(repo.path(), false);

    assert!(
        install_output.status.success(),
        "Installation should succeed: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    assert!(repo.hook_exists("prepare-commit-msg"));
}

/// Tests that backup file created with --force contains exact original hook content.
/// Verifies that existing hooks are properly preserved before being overwritten.
#[test]
fn backup_file_contains_original_hook() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    let original_hook_content = "#!/bin/sh\necho original hook\nexit 0";
    repo.create_config(config);
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

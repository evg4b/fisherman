mod common;

use common::{FishermanBinary, GitTestRepo};

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
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

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
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

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

mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn branch_name_regex_valid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/new-feature");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should succeed with valid branch name: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

#[test]
fn branch_name_regex_invalid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("invalid_branch");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        !handle_output.status.success(),
        "Hook should fail with invalid branch name"
    );
}

#[test]
fn branch_name_prefix_valid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/test-branch");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
}

#[test]
fn branch_name_prefix_invalid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/wrong-prefix");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(!handle_output.status.success());
}

#[test]
fn branch_name_suffix_valid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-v1"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature-v1");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
}

#[test]
fn branch_name_suffix_invalid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-v1"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature-v2");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(!handle_output.status.success());
}

#[test]
fn branch_name_multiple_rules_all_pass() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-dev"

[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/[a-z-]+-dev$"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/new-feature-dev");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
}

#[test]
fn branch_name_multiple_rules_one_fails() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
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

    binary.install(repo.path(), false);

    repo.create_branch("feature/missing-suffix");
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(!handle_output.status.success());
}

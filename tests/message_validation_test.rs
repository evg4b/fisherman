mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn message_regex_valid_pattern() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|test):\\s.+"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("feat: add new feature");

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success(), "Installation failed");
    assert!(repo.hook_exists("commit-msg"));

    repo.write_commit_msg_file("feat: valid commit message");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(
        handle_output.status.success(),
        "Hook should succeed with valid message: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

#[test]
fn message_regex_invalid_pattern() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|test):\\s.+"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial commit");

    let install_output = binary.install(repo.path(), false);
    assert!(install_output.status.success());

    repo.write_commit_msg_file("invalid commit message");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(
        !handle_output.status.success(),
        "Hook should fail with invalid message"
    );
}

#[test]
fn message_prefix_valid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("feat: add feature");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(handle_output.status.success());
}

#[test]
fn message_prefix_invalid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("fix: wrong prefix");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(!handle_output.status.success());
}

#[test]
fn message_suffix_valid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("commit message [skip ci]");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(handle_output.status.success());
}

#[test]
fn message_suffix_invalid() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("commit message without suffix");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(!handle_output.status.success());
}

#[test]
fn message_multiple_rules_all_pass() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [done]"

[[hooks.commit-msg]]
type = "message-regex"
regex = ".*feature.*"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("feat: add new feature [done]");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(handle_output.status.success());
}

#[test]
fn message_multiple_rules_one_fails() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [done]"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.write_commit_msg_file("feat: missing suffix");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(!handle_output.status.success());
}

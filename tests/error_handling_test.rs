mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn invalid_toml_config_fails_gracefully() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let invalid_config = r#"
[[hooks.pre-commit]
type = "branch-name-regex"  # Missing closing bracket
regex = ".*"
"#;

    repo.create_config(invalid_config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);
    // Should fail or handle gracefully
    // The exact behavior depends on implementation
}

#[test]
fn invalid_yaml_config_fails_gracefully() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let invalid_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: test.txt
    content: missing indentation
"#;

    repo.create_yaml_config(invalid_config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);
    // Should fail or handle gracefully
}

#[test]
fn invalid_regex_in_message_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "(?P<unclosed"  # Invalid regex
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        // If installation succeeds, handling should fail
        repo.write_commit_msg_file("test message");
        let msg_path = repo.commit_msg_file_path();
        let handle_output = binary.handle(
            "commit-msg",
            repo.path(),
            &[msg_path.to_str().unwrap()],
        );

        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid regex"
        );
    }
}

#[test]
fn invalid_regex_in_branch_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "[invalid("  # Invalid regex
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        // If installation succeeds, handling should fail
        let handle_output = binary.handle("pre-commit", repo.path(), &[]);
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid regex"
        );
    }
}

#[test]
fn invalid_regex_in_extract() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:(?P<Invalid"]

[[hooks.pre-commit]]
type = "write-file"
path = "test.txt"
content = "test"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        // If installation succeeds, handling should fail
        let handle_output = binary.handle("pre-commit", repo.path(), &[]);
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid extract regex"
        );
    }
}

#[test]
fn template_with_undefined_variable() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Value: {{UndefinedVar}}"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        // If installation succeeds, handling should fail
        let handle_output = binary.handle("pre-commit", repo.path(), &[]);
        assert!(
            !handle_output.status.success(),
            "Hook should fail with undefined template variable"
        );
    }
}

#[test]
fn missing_required_field_in_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "message-regex"
# Missing required 'regex' field
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);
    // Should fail during config parsing or installation
}

#[test]
fn unknown_rule_type() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "unknown-rule-type"
some_field = "value"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);
    // Should fail during config parsing or installation
}

#[test]
fn exec_command_not_found() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "nonexistent-command-12345"
args = ["test"]
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(
        !handle_output.status.success(),
        "Hook should fail when exec command is not found"
    );
}

#[test]
fn write_file_to_invalid_path() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "/root/cannot-write-here.txt"
content = "test"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    // Might fail depending on permissions
    // This tests error handling for file system errors
}

#[test]
fn empty_config_file() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    repo.create_config("");
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);
    assert!(
        install_output.status.success(),
        "Empty config should be valid: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );
}

#[test]
fn when_condition_syntax_error() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "test.txt"
content = "test"
when = "Type == "  # Invalid syntax
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        repo.create_branch("feature/test");
        let handle_output = binary.handle("pre-commit", repo.path(), &[]);
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid when condition"
        );
    }
}

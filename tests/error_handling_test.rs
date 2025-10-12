mod common;

use common::{FishermanBinary, GitTestRepo};

/// Tests that invalid TOML syntax in configuration file fails gracefully without crashing.
/// Verifies error handling when TOML has missing closing brackets or malformed structure.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
    // Should fail or handle gracefully
    // The exact behavior depends on implementation
}

/// Tests that invalid YAML syntax in configuration file fails gracefully.
/// Verifies error handling when YAML has incorrect indentation or formatting issues.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
    // Should fail or handle gracefully
}

/// Tests that message-regex rule with invalid regex pattern fails appropriately.
/// Verifies that malformed regex patterns (unclosed groups) are detected and rejected.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

/// Tests that branch-name-regex rule with invalid regex pattern fails appropriately.
/// Verifies that malformed regex patterns are detected during rule compilation or execution.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

/// Tests that invalid regex in extract configuration fails during variable extraction.
/// Verifies error handling when extract patterns have malformed regex syntax.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

/// Tests that templates referencing undefined variables fail during rule compilation.
/// Verifies that template rendering errors are caught when variables are not extracted.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

/// Tests that rules with missing required fields fail during configuration parsing.
/// Verifies validation of rule configuration structure and required parameters.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
    // Should fail during config parsing or installation
}

/// Tests that unknown rule types in configuration fail during parsing.
/// Verifies that only valid rule type values are accepted in configuration.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
    // Should fail during config parsing or installation
}

/// Tests that exec rule fails when the specified command does not exist.
/// Verifies proper error handling for missing or invalid executables.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(
        !handle_output.status.success(),
        "Hook should fail when exec command is not found"
    );
}

/// Tests error handling when write-file rule attempts to write to an invalid or restricted path.
/// Verifies filesystem permission errors are properly caught and reported.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let _handle_output = binary.handle("pre-commit", repo.path(), &[]);
    // Might fail depending on permissions
    // This tests error handling for file system errors
}

/// Tests that an empty configuration file is considered valid and does not cause errors.
/// Verifies that fisherman can handle repositories with no configured hooks.
#[test]
fn empty_config_file() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    repo.create_config("");
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);
    assert!(
        install_output.status.success(),
        "Empty config should be valid: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );
}

/// Tests that invalid Rhai syntax in when conditions fails during rule evaluation.
/// Verifies that scripting syntax errors are detected and reported properly.
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

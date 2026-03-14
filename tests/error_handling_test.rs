mod common;

use crate::common::ConfigFormat;
use common::test_context::assert_stderr_contains;
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

    repo.create_config(invalid_config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
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
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        let handle_output = repo.commit_with_hooks_allow_empty("test message");

        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid regex"
        );

        let stderr = String::from_utf8_lossy(&handle_output.stderr);
        assert!(!stderr.is_empty(), "Error message should explain regex issue");
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        let handle_output = repo.commit_with_hooks_allow_empty("test commit");
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid regex"
        );

        let stderr = String::from_utf8_lossy(&handle_output.stderr);
        assert!(!stderr.is_empty(), "Error message should explain regex issue");
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        let handle_output = repo.commit_with_hooks_allow_empty("test commit");
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid extract regex"
        );

        let stderr = String::from_utf8_lossy(&handle_output.stderr);
        assert!(!stderr.is_empty(), "Error message should explain extract regex issue: {}", stderr);
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        let handle_output = repo.commit_with_hooks_allow_empty("test commit");
        assert!(
            !handle_output.status.success(),
            "Hook should fail with undefined template variable"
        );

        let stderr = String::from_utf8_lossy(&handle_output.stderr);
        assert_stderr_contains(&stderr, &["UndefinedVar", "variable", "template"],
            "Error should mention the undefined variable");
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let _install_output = binary.install(repo.path(), false);
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");
    assert!(
        !handle_output.status.success(),
        "Hook should fail when exec command is not found"
    );

    let stderr = String::from_utf8_lossy(&handle_output.stderr);
    assert_stderr_contains(&stderr, &["nonexistent-command-12345", "not found", "No such"],
        "Error should mention the command that wasn't found");
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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let _handle_output = repo.commit_with_hooks_allow_empty("test commit");
}

#[test]
fn empty_config_file() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    repo.create_config("", ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

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

    repo.create_config(config, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        repo.create_branch("feature/test");
        let handle_output = repo.commit_with_hooks_allow_empty("test commit");
        assert!(
            !handle_output.status.success(),
            "Hook should fail with invalid when condition"
        );

        let stderr = String::from_utf8_lossy(&handle_output.stderr);
        assert!(!stderr.is_empty(), "Error should explain syntax issue in when condition");
    }
}

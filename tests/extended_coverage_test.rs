mod common;

use common::test_context::TestContext;

// Unicode handling tests

/// Tests that commit messages with Unicode characters are properly validated.
/// Verifies support for international characters and emojis in message-regex rules.
#[test]
fn unicode_in_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_success("feat: Add 日本語 support with émojis 🎉");
}

/// Tests that branch names containing Unicode characters are correctly validated.
/// Verifies that international characters in branch names work with regex validation.
#[test]
fn unicode_in_branch_name() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/support-日本語");
    ctx.git_commit_allow_empty_success("test commit");
}

/// Tests that template variables can contain and render Unicode characters correctly.
/// Verifies that extraction and template rendering preserve international characters.
#[test]
fn unicode_in_template_variable() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Name>.+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "branch-name.txt"
content = "Branch: {{Name}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/日本語-support");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("branch-name.txt");
    assert!(content.contains("日本語"));
}

// prepare-commit-msg hook tests

/// Tests that prepare-commit-msg hook executes correctly with basic write-file rule.
/// Verifies that this hook type is properly supported and receives correct arguments.
/// NOTE: prepare-commit-msg is called directly because Git triggers it before user edits the message,
/// making it difficult to test through natural Git commands in an automated test environment.
#[test]
fn prepare_commit_msg_hook_execution() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.prepare-commit-msg]]
type = "write-file"
path = "prepare-executed.txt"
content = "prepare-commit-msg ran"
"#;

    ctx.setup_and_install(config);

    let msg_path = ctx.repo.commit_msg_file_path();
    let handle_output = ctx.binary.handle(
        "prepare-commit-msg",
        ctx.repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(
        handle_output.status.success(),
        "prepare-commit-msg should execute: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
    assert!(ctx.repo.file_exists("prepare-executed.txt"));
}

/// Tests that template variables are correctly rendered in prepare-commit-msg hook.
/// Verifies variable extraction and template substitution work in this hook context.
/// NOTE: prepare-commit-msg is called directly because Git triggers it before user edits the message,
/// making it difficult to test through natural Git commands in an automated test environment.
#[test]
fn prepare_commit_msg_with_template_variable() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.prepare-commit-msg]]
type = "write-file"
path = "commit-template.txt"
content = "{{Type}}: [{{Ticket}}] "
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/PROJ-456");

    let msg_path = ctx.repo.commit_msg_file_path();
    let handle_output = ctx.binary.handle(
        "prepare-commit-msg",
        ctx.repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(handle_output.status.success());
    assert!(ctx.repo.file_exists("commit-template.txt"));
    let content = ctx.repo.read_file("commit-template.txt");
    assert_eq!(content, "feature: [PROJ-456] ");
}

// Edge cases in conditionals

/// Tests that when condition referencing undefined variable causes rule compilation to fail.
/// Verifies that undefined variables in conditionals are properly detected as errors.
#[test]
fn conditional_with_undefined_variable_fails() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "test"
when = "UndefinedVar == \"value\""
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_failure("test commit");
}

/// Tests that is_def_var() function returns true when variable is extracted and defined.
/// Verifies conditional rule execution based on variable presence using is_def_var().
#[test]
fn conditional_with_is_def_var_true() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Feature is defined"
when = "is_def_var(\"Feature\")"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/auth");
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("output.txt"), "Feature is defined");
}

/// Tests that is_def_var() returns false when optional variable is not extracted.
/// Verifies conditional branching with negative is_def_var() checks for optional extracts.
#[test]
fn conditional_with_is_def_var_false() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Feature is defined"
when = "is_def_var(\"Feature\")"

[[hooks.pre-commit]]
type = "write-file"
path = "fallback.txt"
content = "Feature not defined"
when = "!is_def_var(\"Feature\")"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/test");
    ctx.git_commit_allow_empty_success("test commit");

    assert!(!ctx.repo.file_exists("output.txt"));
    assert!(ctx.repo.file_exists("fallback.txt"));
}

// Environment variable tests

/// Tests that shell scripts can access multiple custom environment variables simultaneously.
/// Verifies that complex env configurations are properly passed to shell context.
#[test]
fn shell_script_with_multiple_env_vars() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "if \"%VAR1%\" == \"value1\" if \"%VAR2%\" == \"value2\" exit 0"
env = { VAR1 = "value1", VAR2 = "value2" }
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = """
#!/bin/sh
if [ "$VAR1" = "value1" ] && [ "$VAR2" = "value2" ]; then
    exit 0
else
    exit 1
fi
"""
env = { VAR1 = "value1", VAR2 = "value2" }
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_success("test commit");
}

/// Tests that environment variables can use template substitution from extracted variables.
/// Verifies template rendering works in env variable values for exec commands.
#[test]
fn exec_with_templated_env_vars() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "%FEATURE_NAME%"]
env = { FEATURE_NAME = "{{Feature}}" }
"#;

    #[cfg(not(windows))]
    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "exec"
command = "sh"
args = ["-c", "test \"$FEATURE_NAME\" = \"payment\""]
env = { FEATURE_NAME = "{{Feature}}" }
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/payment");
    ctx.git_commit_allow_empty_success("test commit");
}

// Edge cases

/// Tests that message validation correctly rejects empty commit messages.
/// Verifies that regex patterns requiring content fail for empty strings.
#[test]
fn empty_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_failure("");
}

/// Tests that very long commit messages (10000+ characters) are handled correctly.
/// Verifies that message validation works without performance or memory issues on large inputs.
#[test]
fn very_long_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    let long_message = "a".repeat(10000);
    ctx.git_commit_allow_empty_success(&long_message);
}

/// Tests that commit messages containing only whitespace characters are rejected by Git.
/// Git itself prevents whitespace-only commit messages, treating them as empty.
/// This test verifies that Git's built-in validation works as expected.
#[test]
fn whitespace_only_commit_message_rejected_by_git() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = ".*"
"#;

    ctx.setup_and_install(config);
    // Git rejects whitespace-only messages as empty, so this should fail
    let output = ctx.git_commit_allow_empty("   \n   \t   ");
    assert!(!output.status.success(), "Git should reject whitespace-only commit messages");
}

/// Tests that write-file rule preserves special characters in content without escaping.
/// Verifies that shell metacharacters, quotes, and symbols are written literally.
#[test]
fn write_file_with_special_characters_in_content() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "special.txt"
content = "Line with $VAR and `backticks` and \"quotes\" and 'apostrophes'"
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("special.txt");
    assert!(content.contains("$VAR"));
    assert!(content.contains("`backticks`"));
    assert!(content.contains("\"quotes\""));
}

/// Tests that hierarchical branch names with multiple slashes are extracted correctly.
/// Verifies complex regex patterns can capture multiple path segments from branch names.
#[test]
fn branch_name_with_slashes() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^(?P<Category>[^/]+)/(?P<Subcategory>[^/]+)/(?P<Name>.+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "hierarchy.txt"
content = "{{Category}}/{{Subcategory}}/{{Name}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/ui/button-component");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("hierarchy.txt");
    assert_eq!(content, "feature/ui/button-component");
}

/// Tests that multiple synchronous branch validation rules all execute and pass together.
/// Verifies that prefix, regex, and suffix rules can be combined successfully.
#[test]
fn multiple_rules_with_mixed_success_sync() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/[a-z-]+$"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-ready"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/new-feature-ready");
    ctx.git_commit_allow_empty_success("test commit");
}

/// Tests that regex patterns with escaped characters (backslashes, digits) work correctly.
/// Verifies proper handling of escape sequences in extraction patterns and template rendering.
#[test]
fn regex_with_escaped_characters() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Ticket>[A-Z]+-\\d+)-(?P<Priority>high|low)"]

[[hooks.pre-commit]]
type = "write-file"
path = "ticket.txt"
content = "{{Ticket}} - {{Priority}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/PROJ-123-high");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("ticket.txt");
    assert_eq!(content, "PROJ-123 - high");
}

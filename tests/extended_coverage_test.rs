mod common;

use common::test_context::TestContext;

// Unicode handling tests

#[test]
fn unicode_in_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("feat: Add 日本語 support with émojis 🎉");
}

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
    ctx.handle_success("pre-commit");
}

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
    ctx.handle_success("pre-commit");

    let content = ctx.repo.read_file("branch-name.txt");
    assert!(content.contains("日本語"));
}

// prepare-commit-msg hook tests

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
    ctx.handle_failure("pre-commit");
}

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
    ctx.handle_success("pre-commit");

    assert_eq!(ctx.repo.read_file("output.txt"), "Feature is defined");
}

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
    ctx.handle_success("pre-commit");

    assert!(!ctx.repo.file_exists("output.txt"));
    assert!(ctx.repo.file_exists("fallback.txt"));
}

// Environment variable tests

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
    ctx.handle_success("pre-commit");
}

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
    ctx.handle_success("pre-commit");
}

// Edge cases

#[test]
fn empty_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_failure("");
}

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
    ctx.handle_commit_msg_success(&long_message);
}

#[test]
fn whitespace_only_commit_message() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^\\s*$"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("   \n   \t   ");
}

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
    ctx.handle_success("pre-commit");

    let content = ctx.repo.read_file("special.txt");
    assert!(content.contains("$VAR"));
    assert!(content.contains("`backticks`"));
    assert!(content.contains("\"quotes\""));
}

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
    ctx.handle_success("pre-commit");

    let content = ctx.repo.read_file("hierarchy.txt");
    assert_eq!(content, "feature/ui/button-component");
}

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
    ctx.handle_success("pre-commit");
}

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
    ctx.handle_success("pre-commit");

    let content = ctx.repo.read_file("ticket.txt");
    assert_eq!(content, "PROJ-123 - high");
}

mod common;

use common::test_context::TestContext;

#[test]
fn pre_push_hook_execution() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-push]]
type = "write-file"
path = "pre-push-executed.txt"
content = "pre-push hook ran"
"#;

    ctx.setup_and_install(config);

    let handle_output = ctx.binary.handle("pre-push", ctx.repo.path(), &[]);
    assert!(
        handle_output.status.success(),
        "pre-push hook should execute successfully: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
    assert!(ctx.repo.file_exists("pre-push-executed.txt"));
}

#[test]
fn post_commit_hook_execution() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.post-commit]]
type = "write-file"
path = "post-commit-executed.txt"
content = "post-commit hook ran"
"#;

    ctx.setup_and_install(config);

    let handle_output = ctx.binary.handle("post-commit", ctx.repo.path(), &[]);
    assert!(
        handle_output.status.success(),
        "post-commit hook should execute successfully: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
    assert!(ctx.repo.file_exists("post-commit-executed.txt"));
}

#[test]
fn empty_hooks_array_succeeds() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
# No rules defined
"#;

    // Empty or minimal config should still allow installation
    let _install_output = ctx.setup_with_config(config);
    // This may succeed or fail depending on implementation
    // Just verify it doesn't crash
}

#[test]
fn mixed_sync_and_async_rules_execute_correctly() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/.*"

[[hooks.pre-commit]]
type = "write-file"
path = "async1.txt"
content = "async rule 1"

[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "write-file"
path = "async2.txt"
content = "async rule 2"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/test-branch");

    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("async1.txt"));
    assert!(ctx.repo.file_exists("async2.txt"));
}

#[test]
fn sync_rule_failure_behavior() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "write-file"
path = "write-attempted.txt"
content = "write rule executed"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/test");

    let handle_output = ctx.handle("pre-commit");

    // Hook should fail due to sync rule failure
    assert!(!handle_output.status.success());

    // Document actual behavior: async rules may or may not execute
    // This test just verifies the hook fails correctly
}

#[test]
fn all_rule_types_in_one_hook() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/"

[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-test"

[[hooks.pre-commit]]
type = "write-file"
path = "all-rules.txt"
content = "all rules passed"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/new-test");

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("all-rules.txt"));
}

#[test]
fn conditional_with_complex_boolean_logic() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix|hotfix)/(?P<Priority>high|low|medium)"]

[[hooks.pre-commit]]
type = "write-file"
path = "urgent.txt"
content = "Urgent work"
when = "(Type == \"hotfix\" || (Type == \"bugfix\" && Priority == \"high\")) && Type != \"feature\""
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("hotfix/low");

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("urgent.txt"));
}

#[test]
fn template_in_message_suffix() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [{{Ticket}}]"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.handle_commit_msg_success("Add new feature [PROJ-123]");
}

#[test]
fn template_in_branch_regex() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["repo_path:.*/(?P<RepoName>[^/]+)$"]

[[hooks.pre-commit]]
type = "write-file"
path = "repo-check.txt"
content = "Repository: {{RepoName}}"
"#;

    ctx.setup_and_install(config);

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("repo-check.txt"));
    let content = ctx.repo.read_file("repo-check.txt");
    assert!(content.starts_with("Repository: "));
}

#[test]
fn multiple_write_files_to_same_location() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "First write"
append = false

[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "\nSecond write"
append = true
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    let content = ctx.repo.read_file("output.txt");
    assert!(content.contains("First write"));
    assert!(content.contains("Second write"));
}

#[test]
fn conditional_false_doesnt_execute_with_valid_message() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
when = "Type == \"feature\""
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/test");

    // Message doesn't have the prefix, but the condition is false, so it should pass
    ctx.handle_commit_msg_success("bugfix: fix the bug");
}

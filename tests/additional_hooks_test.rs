mod common;

use common::test_context::TestContext;

/// Tests that post-merge hook executes successfully.
/// Verifies that post-merge hooks are properly supported and execute write-file rules.
/// NOTE: post-merge is called directly because it requires a complex merge operation setup
/// to test through natural Git commands in test environment.
#[test]
fn post_merge_hook_execution() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.post-merge]]
type = "write-file"
path = "merge-executed.txt"
content = "post-merge ran"
"#;

    ctx.setup_and_install(config);
    let output = ctx.binary.handle("post-merge", ctx.repo.path(), &[]);

    assert!(
        output.status.success(),
        "post-merge hook should execute successfully: {}",
        String::from_utf8_lossy(&output.stderr)
    );
    assert!(ctx.repo.file_exists("merge-executed.txt"));
}

/// Tests that post-checkout hook executes successfully.
/// Verifies that post-checkout hooks are properly supported.
#[test]
fn post_checkout_hook_execution() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.post-checkout]]
type = "write-file"
path = "checkout-executed.txt"
content = "post-checkout ran"
"#;

    ctx.setup_and_install(config);

    // Create and checkout a new branch - this triggers post-checkout hook
    ctx.git_checkout_new_branch("test-branch");

    assert!(
        ctx.repo.file_exists("checkout-executed.txt"),
        "post-checkout hook should create file"
    );
}

/// Tests that pre-receive hook executes successfully.
/// Verifies that pre-receive hooks are properly supported.
/// NOTE: pre-receive is a server-side hook called directly because it's triggered during
/// git push on the remote server, which is complex to test in local test environment.
#[test]
fn pre_receive_hook_execution() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.pre-receive]]
type = "write-file"
path = "receive-executed.txt"
content = "pre-receive ran"
"#;

    ctx.setup_and_install(config);
    let output = ctx.binary.handle("pre-receive", ctx.repo.path(), &[]);

    assert!(
        output.status.success(),
        "pre-receive hook should execute successfully: {}",
        String::from_utf8_lossy(&output.stderr)
    );
    assert!(ctx.repo.file_exists("receive-executed.txt"));
}

/// Tests that very long branch names are handled correctly.
/// Verifies that branch name validation works without length limits.
#[test]
fn very_long_branch_name() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    ctx.setup_and_install(config);

    // Create a branch name with 200 characters
    let long_name = format!("feature/{}", "a".repeat(192));
    ctx.repo.create_branch(&long_name);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(
        output.status.success(),
        "Should handle very long branch names: {}",
        String::from_utf8_lossy(&output.stderr)
    );
}

/// Tests that branch names with special characters work correctly.
/// Verifies regex matching with dots, underscores, and other valid Git characters.
#[test]
fn branch_name_with_special_characters() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/[a-z0-9._-]+$"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/test_feature.v1.2-beta");

    ctx.git_commit_allow_empty_success("test commit");
}

/// Tests that write-file with append mode creates file if it doesn't exist.
/// Verifies that append mode works correctly even when target file is missing.
#[test]
fn write_file_append_to_nonexistent() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "new-file.txt"
content = "content"
append = true
"#;

    ctx.setup_and_install(config);
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("new-file.txt"), "content");
}

/// Tests combining conditional execution with template variables.
/// Verifies complex conditional logic with multiple variables.
#[test]
fn conditional_with_multiple_template_variables() {
    let ctx = TestContext::new();

    let config = r#"
extract = [
    "branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)",
    "repo_path:.*/(?P<RepoName>[^/]+)$"
]

[[hooks.pre-commit]]
type = "write-file"
path = "context.txt"
content = "{{Type}}: {{Ticket}} in {{RepoName}}"
when = "Type == \"feature\" && is_def_var(\"Ticket\")"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/PROJ-789");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("context.txt");
    assert!(content.contains("feature: PROJ-789"));
}

/// Tests that explain command works for hooks with no configured rules.
/// Verifies graceful handling when explaining unconfigured hooks.
#[test]
fn explain_unconfigured_hook() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    ctx.setup_and_install(config);

    // Explain a different hook that has no rules
    let output = ctx.binary.explain("pre-push", ctx.repo.path());

    assert!(
        output.status.success(),
        "Explain should succeed even for unconfigured hooks"
    );
}

/// Tests that multiple extraction patterns can extract from the same source.
/// Verifies that branch patterns don't conflict when extracting different groups.
#[test]
fn multiple_extractions_same_source() {
    let ctx = TestContext::new();

    let config = r#"
extract = [
    "branch:^(?P<Type>feature|bugfix)",
    "branch:^[^/]+/(?P<Name>.+)"
]

[[hooks.pre-commit]]
type = "write-file"
path = "extracted.txt"
content = "{{Type}}: {{Name}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/auth-system");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("extracted.txt");
    assert_eq!(content, "feature: auth-system");
}

/// Tests that commit messages with newlines are validated correctly.
/// Verifies multiline commit message handling in message-regex rules.
#[test]
fn commit_message_with_newlines() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^feat: .+"
"#;

    ctx.setup_and_install(config);

    let multiline_msg = "feat: add new feature\n\nThis is a longer description\nwith multiple lines";
    ctx.git_commit_allow_empty_success(multiline_msg);
}

/// Tests that optional extraction doesn't fail when pattern doesn't match.
/// Verifies that hooks succeed when optional variables aren't extracted.
#[test]
fn optional_extraction_no_match() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "executed.txt"
content = "Hook executed"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/issue");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("executed.txt");
    assert_eq!(content, "Hook executed");
}

mod common;

use common::test_context::TestContext;

#[test]
#[cfg(feature = "integration-tests")]
fn post_merge_hook_execution() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.post-merge]]
type = "write-file"
path = "merge-executed.txt"
content = "post-merge ran"
"#;

    ctx.setup_and_install_old(config);

    ctx.repo.create_file("file1.txt", "content1");
    ctx.repo.commit("initial commit");

    ctx.repo.create_branch("feature-branch");
    ctx.repo.create_file("file2.txt", "content2");
    ctx.repo.commit("feature commit");

    ctx.repo.checkout("master");
    let output = ctx.repo.git(&["merge", "feature-branch", "--no-edit"]);

    assert!(
        output.status.success(),
        "Merge should succeed: {}",
        String::from_utf8_lossy(&output.stderr)
    );
    assert!(ctx.repo.file_exists("merge-executed.txt"), "post-merge hook should have created file");
}

#[test]
#[cfg(feature = "integration-tests")]
fn post_checkout_hook_execution() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.post-checkout]]
type = "write-file"
path = "checkout-executed.txt"
content = "post-checkout ran"
"#;

    ctx.setup_and_install_old(config);

    ctx.git_checkout_new_branch("test-branch");

    assert!(
        ctx.repo.file_exists("checkout-executed.txt"),
        "post-checkout hook should create file"
    );
}

// NOTE: pre-receive is a server-side hook that runs during git push on the remote repository.

#[test]
#[cfg(feature = "integration-tests")]
fn very_long_branch_name() {
    let ctx = TestContext::new();
    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    ctx.setup_and_install_old(config);

    let long_name = format!("feature/{}", "a".repeat(192));
    ctx.repo.create_branch(&long_name);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(
        output.status.success(),
        "Should handle very long branch names: {}",
        String::from_utf8_lossy(&output.stderr)
    );
}

#[test]
#[cfg(feature = "integration-tests")]
fn branch_name_with_special_characters() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/[a-z0-9._-]+$"
"#;

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/test_feature.v1.2-beta");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
#[cfg(feature = "integration-tests")]
fn write_file_append_to_nonexistent() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "new-file.txt"
content = "content"
append = true
"#;

    ctx.setup_and_install_old(config);
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("new-file.txt"), "content");
}

#[test]
#[cfg(feature = "integration-tests")]
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

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/PROJ-789");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("context.txt");
    assert!(content.contains("feature: PROJ-789"));
}

#[test]
#[cfg(feature = "integration-tests")]
fn explain_unconfigured_hook() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = ".*"
"#;

    ctx.setup_and_install_old(config);

    let output = ctx.binary.explain("pre-push", ctx.repo.path());

    assert!(
        output.status.success(),
        "Explain should succeed even for unconfigured hooks"
    );
}

#[test]
#[cfg(feature = "integration-tests")]
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

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/auth-system");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("extracted.txt");
    assert_eq!(content, "feature: auth-system");
}

#[test]
#[cfg(feature = "integration-tests")]
fn commit_message_with_newlines() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^feat: .+"
"#;

    ctx.setup_and_install_old(config);

    let multiline_msg = "feat: add new feature\n\nThis is a longer description\nwith multiple lines";
    ctx.git_commit_allow_empty_success(multiline_msg);
}

#[test]
#[cfg(feature = "integration-tests")]
fn optional_extraction_no_match() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "executed.txt"
content = "Hook executed"
"#;

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("bugfix/issue");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("executed.txt");
    assert_eq!(content, "Hook executed");
}

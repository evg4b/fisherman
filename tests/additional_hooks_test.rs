mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use fisherman_core::{t, BranchNameRegexRule, Configuration, Expression, GitHook, WriteFileRule};

#[test]
fn post_merge_hook_execution() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PostMerge => [
            rule!(WriteFileRule {
                path: "merge-executed.txt".into(),
                content: "post-merge ran".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

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
fn post_checkout_hook_execution() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PostCheckout => [
            rule!(WriteFileRule {
                path: "checkout-executed.txt".into(),
                content: "post-checkout ran".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    ctx.git_checkout_new_branch("test-branch");

    assert!(
        ctx.repo.file_exists("checkout-executed.txt"),
        "post-checkout hook should create file"
    );
}

#[test]
fn branch_name_with_special_characters() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^feature/[a-z0-9._-]+$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/test_feature.v1.2-beta");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn write_file_append_to_nonexistent() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "new-file.txt".into(),
                content: "content".into(),
                append: Some(true),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("new-file.txt"), "content");
}

#[test]
fn conditional_with_multiple_template_variables() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "context.txt".into(),
                content: t!("{{Type}}: {{Ticket}} in {{RepoName}}"),
                append: None,
            }, when = Expression::new("Type == \"feature\" && is_def_var(\"Ticket\")"))
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"),
            String::from("repo_path:.*[\\/\\\\](?P<RepoName>[^\\/\\\\]+)$"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-789");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("context.txt");
    assert!(content.contains("feature: PROJ-789"));
}

#[test]
fn explain_unconfigured_hook() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: ".*".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let output = ctx.binary.explain("pre-push", ctx.repo.path());

    assert!(
        output.status.success(),
        "Explain should succeed even for unconfigured hooks"
    );
}

#[test]
fn multiple_extractions_same_source() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "extracted.txt".into(),
                content: t!("{{Type}}: {{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
            String::from("branch:^[^/]+/(?P<Name>.+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth-system");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("extracted.txt");
    assert_eq!(content, "feature: auth-system");
}

#[test]
fn commit_message_with_newlines() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(fisherman_core::CommitMessageRegexRule {
                when: None,
                expression: "^feat: .+".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let multiline_msg = "feat: add new feature\n\nThis is a longer description\nwith multiple lines";
    ctx.git_commit_allow_empty_success(multiline_msg);
}

#[test]
fn optional_extraction_no_match() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "executed.txt".into(),
                content: "Hook executed".into(),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch?:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/issue");

    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("executed.txt");
    assert_eq!(content, "Hook executed");
}

mod common;

use crate::common::configuration::serialize_configuration;
use crate::common::ConfigFormat;
use common::test_context::TestContext;
use fisherman_core::{t, BranchNamePrefixRule, BranchNameRegexRule, BranchNameSuffixRule, CommitMessageSuffixRule};
use fisherman_core::{Configuration, Expression, GitHook, WriteFileRule};

#[test]
fn post_commit_hook_execution() {
    let ctx = TestContext::new();

    let path = "post-commit-executed.txt";
    let content = "post-commit hook ran";

    let config = config!(
        GitHook::PostCommit => [
            rule!(WriteFileRule {
                path: path.into(),
                content: content.into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists(path));
    assert_eq!(ctx.repo.read_file(path), "post-commit hook ran");
}

#[test]
fn empty_hooks_array_succeeds() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => []
    );

    let _install_output = ctx.setup_with_config(
        serialize_configuration(&config, ConfigFormat::Toml).as_str()
    );
}

#[test]
fn mixed_sync_and_async_rules_execute_correctly() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^feature/.*".into(),
            }),
            rule!(WriteFileRule {
                path: "async1.txt".into(),
                content: "async rule 1".into(),
                append: None,
            }),
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(WriteFileRule {
                path: "async2.txt".into(),
                content: "async rule 2".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/test-branch");

    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("async1.txt"));
    assert!(ctx.repo.file_exists("async2.txt"));
}

#[test]
fn sync_rule_failure_behavior() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^feature/.*".into(),
            }),
            rule!(WriteFileRule {
                path: "async1.txt".into(),
                content: "async rule 1".into(),
                append: None,
            }),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/test");

    let handle_output = ctx.git_commit_allow_empty("test commit");

    assert!(!handle_output.status.success());
}

#[test]
fn all_rule_types_in_one_hook() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^feature/".into(),
            }),
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(BranchNameSuffixRule {
                suffix: "-test".into(),
            }),
            rule!(WriteFileRule {
                path: "all-rules.txt".into(),
                content: "all rules passed".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/new-test");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("all-rules.txt"));
}

#[test]
fn conditional_with_complex_boolean_logic() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "urgent.txt".into(),
                content: "Urgent work".into(),
                append: None,
            }, when = Expression::new("(Type == \"hotfix\" || (Type == \"bugfix\" && Priority == \"high\")) && Type != \"feature\""))
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix|hotfix)/(?P<Priority>high|low|medium)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("hotfix/low");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("urgent.txt"));
}

#[test]
fn template_in_message_suffix() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageSuffixRule {
                suffix: t!("[{{Ticket}}]")
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Ticket>[A-Z]+-\\d+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.git_commit_allow_empty_success("Add new feature [PROJ-123]");
}

#[test]
fn template_in_branch_regex() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "repo-check.txt".into(),
                content: t!("Repository: {{RepoName}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("repo_path:.*[\\/\\\\](?P<RepoName>[^\\/\\\\]+)$"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("repo-check.txt"));
    let content = ctx.repo.read_file("repo-check.txt");
    assert!(content.starts_with("Repository: "));
}

#[test]
fn conditional_false_doesnt_execute_with_valid_message() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(fisherman_core::CommitMessagePrefixRule {
                prefix: "feat: ".into(),
            }, when = Expression::new("Type == \"feature\""))
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/test");

    ctx.git_commit_allow_empty_success("bugfix: fix the bug");
}

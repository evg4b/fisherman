mod common;

use crate::common::configuration::serialize_configuration;
use crate::common::ConfigFormat;
use common::test_context::TestContext;
use core::configuration::Configuration;
use core::hooks::GitHook;

#[test]
fn post_commit_hook_execution() {
    let ctx = TestContext::new();

    let path = "post-commit-executed.txt";
    let content = "post-commit hook ran";

    let config = config!(
        GitHook::PostCommit => [
            rule!(WriteFileRule {
                path: String::from(path),
                content: String::from(content),
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
            rule!(RuleParams::BranchNameRegex {
                regex: String::from("^feature/.*"),
            }),
            rule!(WriteFileRule {
                path: String::from("async1.txt"),
                content: String::from("async rule 1"),
                append: None,
            }),
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("feature/"),
            }),
            rule!(WriteFileRule {
                path: String::from("async2.txt"),
                content: String::from("async rule 2"),
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
            rule!(RuleParams::BranchNameRegex {
                regex: String::from("^feature/.*"),
            }),
            rule!(WriteFileRule {
                path: String::from("async1.txt"),
                content: String::from("async rule 1"),
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

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/new-test");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("all-rules.txt"));
}

#[test]
fn conditional_with_complex_boolean_logic() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("urgent.txt"),
                    content: String::from("Urgent work"),
                    append: None,
                },
                when = String::from("(Type == \"hotfix\" || (Type == \"bugfix\" && Priority == \"high\")) && Type != \"feature\"")
            )
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

    let config = r#"
extract = ["branch:^feature/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [{{Ticket}}]"
"#;

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.git_commit_allow_empty_success("Add new feature [PROJ-123]");
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

    ctx.setup_and_install_old(config);

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("repo-check.txt"));
    let content = ctx.repo.read_file("repo-check.txt");
    assert!(content.starts_with("Repository: "));
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

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("bugfix/test");

    ctx.git_commit_allow_empty_success("bugfix: fix the bug");
}

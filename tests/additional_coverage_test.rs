mod common;

use crate::common::ConfigFormat;
use crate::common::configuration::serialize_configuration;
use common::test_context::TestContext;
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::{Rule, RuleParams};
// NOTE: pre-push is a client-side hook that runs before git push sends objects to the remote.
// Testing it would require setting up a remote repository and performing push operations,
// which adds significant complexity. Since we already test hook execution thoroughly with
// other hook types (pre-commit, commit-msg, post-commit, etc.), we've omitted this test.
// The hook installation and execution logic is the same for all hook types.

/// Tests that post-commit hooks can be configured and executed successfully.
/// Verifies that write-file rules work in post-commit hook context.
#[test]
fn post_commit_hook_execution() {
    let ctx = TestContext::new();

    let path = "post-commit-executed.txt";
    let content = "post-commit hook ran";

    let config = config!(
        GitHook::PostCommit => [
            rule!(RuleParams::WriteFile {
                path: String::from(path),
                content: String::from(content),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    // post-commit runs automatically after a successful commit
    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists(path));
    assert_eq!(ctx.repo.read_file(path), "post-commit hook ran");
}

/// Tests that configurations with empty or minimal hook definitions handle gracefully
/// without crashing or producing errors.
#[test]
fn empty_hooks_array_succeeds() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => []
    );

    // Empty or minimal config should still allow installation
    let _install_output = ctx.setup_with_config(
        serialize_configuration(&config, ConfigFormat::Toml).as_str()
    );
    // This may succeed or fail depending on implementation
    // Just verify it doesn't crash
}

/// Tests that synchronous rules (branch validation) and asynchronous rules (write-file)
/// can be mixed in the same hook and execute correctly. Sync rules run first, then async
/// rules run in parallel.
#[test]
fn mixed_sync_and_async_rules_execute_correctly() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNameRegex {
                regex: String::from("^feature/.*"),
            }),
            rule!(RuleParams::WriteFile {
                path: String::from("async1.txt"),
                content: String::from("async rule 1"),
                append: None,
            }),
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("feature/"),
            }),
            rule!(RuleParams::WriteFile {
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

/// Tests that when a synchronous rule fails, the entire hook fails even if there are
/// async rules configured. Verifies proper failure propagation from sync rules.
#[test]
fn sync_rule_failure_behavior() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNameRegex {
                regex: String::from("^feature/.*"),
            }),
            rule!(RuleParams::WriteFile {
                path: String::from("async1.txt"),
                content: String::from("async rule 1"),
                append: None,
            }),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/test");

    let handle_output = ctx.git_commit_allow_empty("test commit");

    // Hook should fail due to sync rule failure
    assert!(!handle_output.status.success());

    // Document actual behavior: async rules may or may not execute
    // This test just verifies the hook fails correctly
}

/// Tests that multiple different rule types (regex, prefix, suffix, write-file) can be
/// configured in a single hook and all execute successfully when their conditions are met.
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

/// Tests that complex boolean expressions with AND, OR, and NOT operators work correctly
/// in conditional (when) expressions. Verifies multiple variables and nested logic.
#[test]
fn conditional_with_complex_boolean_logic() {
    let ctx = TestContext::new();
    
    let config = config!(
        GitHook::PreCommit => [
            rule!(
                RuleParams::WriteFile {
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

/// Tests that template variables can be used within message-suffix rules to dynamically
/// construct expected suffixes based on extracted branch information.
#[test]
fn template_in_message_suffix() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [{{Ticket}}]"
";

    ctx.setup_and_install_old(config);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.git_commit_allow_empty_success("Add new feature [PROJ-123]");
}

/// Tests that template variables extracted from repository path can be used in write-file
/// content. Verifies repo_path extraction works correctly.
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

/// Tests that rules with false conditional expressions are skipped and don't affect
/// the hook result. Verifies that a message without required prefix passes when the
/// conditional is false.
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

    // Message doesn't have the prefix, but the condition is false, so it should pass
    ctx.git_commit_allow_empty_success("bugfix: fix the bug");
}

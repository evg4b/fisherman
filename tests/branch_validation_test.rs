mod common;

use common::test_context::{assert_stderr_contains, TestContext};
use common::ConfigFormat;
use fisherman_core::BranchNamePrefixRule;
use fisherman_core::BranchNameRegexRule;
use fisherman_core::BranchNameSuffixRule;
use fisherman_core::Configuration;
use fisherman_core::GitHook;

#[test]
fn branch_name_regex_valid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^(feature|bugfix|hotfix)/[a-z0-9-]+".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/new-feature");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn branch_name_regex_invalid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^(feature|bugfix|hotfix)/[a-z0-9-]+".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("invalid_branch");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "invalid_branch", "regex"],
                           "Error should explain branch name validation failure");
}

#[test]
fn branch_name_prefix_valid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/test-branch");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn branch_name_prefix_invalid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/wrong-prefix");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "prefix", "feature/"],
                           "Error should explain prefix validation failure");
}

#[test]
fn branch_name_suffix_valid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameSuffixRule {
                suffix: "-v1".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature-v1");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn branch_name_suffix_invalid() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameSuffixRule {
                suffix: "-v1".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature-v2");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "suffix", "-v1"],
                           "Error should explain suffix validation failure");
}

#[test]
fn branch_name_multiple_rules_all_pass() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(BranchNameSuffixRule {
                suffix: "-dev".into(),
            }),
            rule!(BranchNameRegexRule {
                expression: "^feature/[a-z-]+-dev$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/new-feature-dev");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn branch_name_multiple_rules_one_fails() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(BranchNameSuffixRule {
                suffix: "-dev".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/missing-suffix");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain which rule failed");
}

mod common;

use common::test_context::{TestContext, assert_stderr_contains};

/// Tests that branch-name-regex rule passes when branch name matches the specified pattern.
/// Verifies regex validation works with valid branch naming conventions.
#[test]
fn branch_name_regex_valid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_regex!("^(feature|bugfix|hotfix)/[a-z0-9-]+"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("feature/new-feature");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-regex rule fails when branch name doesn't match the pattern.
/// Verifies hook correctly rejects invalid branch names.
#[test]
fn branch_name_regex_invalid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_regex!("^(feature|bugfix|hotfix)/[a-z0-9-]+"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("invalid_branch");

    let output = ctx.handle("pre-commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "invalid_branch", "regex"],
        "Error should explain branch name validation failure");
}

/// Tests that branch-name-prefix rule passes when branch name starts with required prefix.
/// Verifies prefix-based branch naming enforcement works correctly.
#[test]
fn branch_name_prefix_valid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_prefix!("feature/"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("feature/test-branch");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-prefix rule fails when branch name has wrong prefix.
/// Verifies hook rejects branches without the required prefix.
#[test]
fn branch_name_prefix_invalid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_prefix!("feature/"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("bugfix/wrong-prefix");

    let output = ctx.handle("pre-commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "prefix", "feature/"],
        "Error should explain prefix validation failure");
}

/// Tests that branch-name-suffix rule passes when branch name ends with required suffix.
/// Verifies suffix-based branch naming enforcement works correctly.
#[test]
fn branch_name_suffix_valid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_suffix!("-v1"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("feature-v1");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-suffix rule fails when branch name has wrong suffix.
/// Verifies hook rejects branches without the required suffix.
#[test]
fn branch_name_suffix_invalid() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_suffix!("-v1"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("feature-v2");

    let output = ctx.handle("pre-commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["branch", "suffix", "-v1"],
        "Error should explain suffix validation failure");
}

/// Tests that multiple branch validation rules (prefix, suffix, regex) can all pass together
/// when the branch name satisfies all conditions simultaneously.
#[test]
fn branch_name_multiple_rules_all_pass() {
    let ctx = TestContext::new();
    let config = config! {
        hooks: {
            "pre-commit" => [
                branch_prefix!("feature/"),
                branch_suffix!("-dev"),
                branch_regex!("^feature/[a-z-]+-dev$"),
            ]
        }
    };

    ctx.setup_and_install(&config);
    ctx.repo.create_branch("feature/new-feature-dev");
    ctx.handle_success("pre-commit");
}

/// Tests that when multiple branch validation rules are configured, the hook fails if any
/// one rule fails, even if others pass. Verifies all rules must pass.
#[test]
fn branch_name_multiple_rules_one_fails() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-dev"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/missing-suffix");

    let output = ctx.handle("pre-commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain which rule failed");
}

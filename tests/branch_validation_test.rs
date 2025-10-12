mod common;

use common::test_context::TestContext;

/// Tests that branch-name-regex rule passes when branch name matches the specified pattern.
/// Verifies regex validation works with valid branch naming conventions.
#[test]
fn branch_name_regex_valid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/new-feature");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-regex rule fails when branch name doesn't match the pattern.
/// Verifies hook correctly rejects invalid branch names.
#[test]
fn branch_name_regex_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("invalid_branch");
    ctx.handle_failure("pre-commit");
}

/// Tests that branch-name-prefix rule passes when branch name starts with required prefix.
/// Verifies prefix-based branch naming enforcement works correctly.
#[test]
fn branch_name_prefix_valid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/test-branch");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-prefix rule fails when branch name has wrong prefix.
/// Verifies hook rejects branches without the required prefix.
#[test]
fn branch_name_prefix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/wrong-prefix");
    ctx.handle_failure("pre-commit");
}

/// Tests that branch-name-suffix rule passes when branch name ends with required suffix.
/// Verifies suffix-based branch naming enforcement works correctly.
#[test]
fn branch_name_suffix_valid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-v1"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature-v1");
    ctx.handle_success("pre-commit");
}

/// Tests that branch-name-suffix rule fails when branch name has wrong suffix.
/// Verifies hook rejects branches without the required suffix.
#[test]
fn branch_name_suffix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-v1"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature-v2");
    ctx.handle_failure("pre-commit");
}

/// Tests that multiple branch validation rules (prefix, suffix, regex) can all pass together
/// when the branch name satisfies all conditions simultaneously.
#[test]
fn branch_name_multiple_rules_all_pass() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-dev"

[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^feature/[a-z-]+-dev$"
"#;

    ctx.setup_and_install(config);
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
    ctx.handle_failure("pre-commit");
}

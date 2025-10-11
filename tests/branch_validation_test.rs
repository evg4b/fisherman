mod common;

use common::TestContext;

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

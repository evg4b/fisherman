mod common;

use common::test_context::{TestContext, assert_stderr_contains};

/// Tests that message-regex rule passes when commit message matches the specified pattern.
/// Verifies regex validation accepts messages following conventional commit format.
#[test]
fn message_regex_valid_pattern() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|test):\\s.+"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("feat: valid commit message");
}

/// Tests that message-regex rule fails when commit message doesn't match required pattern.
/// Verifies that non-conforming messages are properly rejected by regex validation.
#[test]
fn message_regex_invalid_pattern() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|test):\\s.+"
"#;

    ctx.setup_and_install(config);

    let output = ctx.handle_commit_msg("invalid commit message");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain message validation failure");
}

/// Tests that message-prefix rule passes when message starts with the configured prefix.
/// Verifies simple prefix matching for enforcing commit message conventions.
#[test]
fn message_prefix_valid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("feat: add feature");
}

/// Tests that message-prefix rule fails when message doesn't start with required prefix.
/// Verifies that messages with incorrect or missing prefix are rejected.
#[test]
fn message_prefix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
"#;

    ctx.setup_and_install(config);

    let output = ctx.handle_commit_msg("fix: wrong prefix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["message", "prefix", "feat:"],
        "Error should explain prefix validation failure");
}

/// Tests that message-suffix rule passes when message ends with the configured suffix.
/// Verifies suffix matching for enforcing commit message tags like [skip ci].
#[test]
fn message_suffix_valid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("commit message [skip ci]");
}

/// Tests that message-suffix rule fails when message doesn't end with required suffix.
/// Verifies that messages missing the expected suffix are properly rejected.
#[test]
fn message_suffix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
"#;

    ctx.setup_and_install(config);

    let output = ctx.handle_commit_msg("commit message without suffix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["message", "suffix", "[skip ci]"],
        "Error should explain suffix validation failure");
}

/// Tests that multiple message validation rules all pass when message satisfies all criteria.
/// Verifies that prefix, suffix, and regex rules can be combined successfully.
#[test]
fn message_multiple_rules_all_pass() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [done]"

[[hooks.commit-msg]]
type = "message-regex"
regex = ".*feature.*"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_success("feat: add new feature [done]");
}

/// Tests that hook fails when one of multiple message rules doesn't pass.
/// Verifies that all validation rules must succeed for the hook to pass.
#[test]
fn message_multiple_rules_one_fails() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "

[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [done]"
"#;

    ctx.setup_and_install(config);

    let output = ctx.handle_commit_msg("feat: missing suffix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain which rule failed");
}

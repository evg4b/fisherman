mod common;

use common::TestContext;

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

#[test]
fn message_regex_invalid_pattern() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|test):\\s.+"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_failure("invalid commit message");
}

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

#[test]
fn message_prefix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_failure("fix: wrong prefix");
}

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

#[test]
fn message_suffix_invalid() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
"#;

    ctx.setup_and_install(config);
    ctx.handle_commit_msg_failure("commit message without suffix");
}

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
    ctx.handle_commit_msg_failure("feat: missing suffix");
}

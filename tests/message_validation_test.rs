mod common;

use common::test_context::{assert_stderr_contains, TestContext};
use common::ConfigFormat;
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::RuleParams;

#[test]
fn message_regex_valid_pattern() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessageRegex {
                regex: String::from("^(feat|fix|docs|test):\\s.+"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("feat: valid commit message");
}

#[test]
fn message_regex_invalid_pattern() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessageRegex {
                regex: String::from("^(feat|fix|docs|test):\\s.+"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let output = ctx.git_commit_allow_empty("invalid commit message");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain message validation failure");
}

#[test]
fn message_prefix_valid() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessagePrefix {
                prefix: String::from("feat: "),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("feat: add feature");
}

#[test]
fn message_prefix_invalid() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessagePrefix {
                prefix: String::from("feat: "),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let output = ctx.git_commit_allow_empty("fix: wrong prefix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["message", "prefix", "feat:"],
        "Error should explain prefix validation failure");
}

#[test]
fn message_suffix_valid() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessageSuffix {
                suffix: String::from(" [skip ci]"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("commit message [skip ci]");
}

#[test]
fn message_suffix_invalid() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessageSuffix {
                suffix: String::from(" [skip ci]"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let output = ctx.git_commit_allow_empty("commit message without suffix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert_stderr_contains(&stderr, &["message", "suffix", "[skip ci]"],
        "Error should explain suffix validation failure");
}

#[test]
fn message_multiple_rules_all_pass() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessagePrefix {
                prefix: String::from("feat: "),
            }),
            rule!(RuleParams::CommitMessageSuffix {
                suffix: String::from(" [done]"),
            }),
            rule!(RuleParams::CommitMessageRegex {
                regex: String::from(".*feature.*"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("feat: add new feature [done]");
}

#[test]
fn message_multiple_rules_one_fails() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessagePrefix {
                prefix: String::from("feat: "),
            }),
            rule!(RuleParams::CommitMessageSuffix {
                suffix: String::from(" [done]"),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let output = ctx.git_commit_allow_empty("feat: missing suffix");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain which rule failed");
}

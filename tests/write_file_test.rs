mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use core::Configuration;
use core::GitHook;
use core::WriteFileRule;

#[test]
fn write_file_creates_new_file() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output.txt".into(),
                content: "test content".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("output.txt"), "File should be created");
    assert_eq!(ctx.repo.read_file("output.txt"), "test content");
}

#[test]
fn write_file_overwrites_existing() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output.txt".into(),
                content: "new content".into(),
                append: Some(false),
            })
        ]
    );

    ctx.setup_with_history_and_install(&config, ConfigFormat::Toml, &[("initial", &[
        ("test.txt", "initial"),
        ("output.txt", "old content")
    ])]);

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "new content");
}

#[test]
fn write_file_appends_to_existing() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output.txt".into(),
                content: "\nappended content".into(),
                append: Some(true),
            })
        ]
    );

    ctx.setup_with_history_and_install(&config, ConfigFormat::Toml, &[("initial", &[
        ("test.txt", "initial"),
        ("output.txt", "existing content")
    ])]);

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(
        ctx.repo.read_file("output.txt"),
        "existing content\nappended content"
    );
}

#[test]
fn write_file_simple_path() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "simple.txt".into(),
                content: "simple content".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("simple.txt"));
    assert_eq!(ctx.repo.read_file("simple.txt"), "simple content");
}

#[test]
fn write_file_multiple_files() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output1.txt".into(),
                content: "content 1".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output2.txt".into(),
                content: "content 2".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output3.txt".into(),
                content: "content 3".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("output1.txt"));
    assert!(ctx.repo.file_exists("output2.txt"));
    assert!(ctx.repo.file_exists("output3.txt"));
    assert_eq!(ctx.repo.read_file("output1.txt"), "content 1");
    assert_eq!(ctx.repo.read_file("output2.txt"), "content 2");
    assert_eq!(ctx.repo.read_file("output3.txt"), "content 3");
}

#[test]
fn write_file_multiline_content() {
    let ctx = TestContext::new();
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                when: None,
                extract: None,
                path: "output.txt".into(),
                content: "Line 1\nLine 2\nLine 3".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("output.txt"), "Line 1\nLine 2\nLine 3");
}

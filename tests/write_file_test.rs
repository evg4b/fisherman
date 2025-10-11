mod common;

use common::test_context::TestContext;

#[test]
fn write_file_creates_new_file() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "test content"
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("output.txt"), "File should be created");
    assert_eq!(ctx.repo.read_file("output.txt"), "test content");
}

#[test]
fn write_file_overwrites_existing() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "new content"
append = false
"#;

    ctx.setup_with_history(config, &[("initial", &[
        ("test.txt", "initial"),
        ("output.txt", "old content")
    ])]);

    ctx.handle_success("pre-commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "new content");
}

#[test]
fn write_file_appends_to_existing() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "\nappended content"
append = true
"#;

    ctx.setup_with_history(config, &[("initial", &[
        ("test.txt", "initial"),
        ("output.txt", "existing content")
    ])]);

    ctx.handle_success("pre-commit");
    assert_eq!(
        ctx.repo.read_file("output.txt"),
        "existing content\nappended content"
    );
}

#[test]
fn write_file_simple_path() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "simple.txt"
content = "simple content"
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("simple.txt"));
    assert_eq!(ctx.repo.read_file("simple.txt"), "simple content");
}

#[test]
fn write_file_multiple_files() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output1.txt"
content = "content 1"

[[hooks.pre-commit]]
type = "write-file"
path = "output2.txt"
content = "content 2"

[[hooks.pre-commit]]
type = "write-file"
path = "output3.txt"
content = "content 3"
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

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

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Line 1\nLine 2\nLine 3"
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert_eq!(ctx.repo.read_file("output.txt"), "Line 1\nLine 2\nLine 3");
}

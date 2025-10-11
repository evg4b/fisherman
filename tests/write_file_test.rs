mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn write_file_creates_new_file() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "test content"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should succeed: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
    assert!(repo.file_exists("output.txt"), "File should be created");
    assert_eq!(repo.read_file("output.txt"), "test content");
}

#[test]
fn write_file_overwrites_existing() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "new content"
append = false
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    repo.create_file("output.txt", "old content");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert_eq!(repo.read_file("output.txt"), "new content");
}

#[test]
fn write_file_appends_to_existing() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "\nappended content"
append = true
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    repo.create_file("output.txt", "existing content");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert_eq!(
        repo.read_file("output.txt"),
        "existing content\nappended content"
    );
}

#[test]
fn write_file_simple_path() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "simple.txt"
content = "simple content"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should succeed: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
    assert!(repo.file_exists("simple.txt"));
    assert_eq!(repo.read_file("simple.txt"), "simple content");
}

#[test]
fn write_file_multiple_files() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert!(repo.file_exists("output1.txt"));
    assert!(repo.file_exists("output2.txt"));
    assert!(repo.file_exists("output3.txt"));
    assert_eq!(repo.read_file("output1.txt"), "content 1");
    assert_eq!(repo.read_file("output2.txt"), "content 2");
    assert_eq!(repo.read_file("output3.txt"), "content 3");
}

#[test]
fn write_file_multiline_content() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Line 1\nLine 2\nLine 3"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert_eq!(repo.read_file("output.txt"), "Line 1\nLine 2\nLine 3");
}

mod common;

use common::test_context::TestContext;
use common::ConfigFormat;

#[test]
fn test_suppress_files_rule() {
    let context = TestContext::new();
    let config = r#"
[[hooks.pre-commit]]
type = "suppress-files"
glob = "secret.txt"
"#;
    context.setup_and_install_old(config);

    // Create a normal file and commit - should succeed
    context.repo.create_file("normal.txt", "content");
    context.repo.git(&["add", "normal.txt"]);
    let output = context.repo.commit("commit normal file");
    assert!(output.status.success(), "Commit should succeed for normal file");

    // Create a suppressed file and try to commit - should fail
    context.repo.create_file("secret.txt", "confidential");
    context.repo.git(&["add", "secret.txt"]);
    let output = context.repo.commit("commit secret file");
    assert!(!output.status.success(), "Commit should fail for suppressed file");
    
    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(stderr.contains("The following files are suppressed from being committed: secret.txt"));
}

#[test]
fn test_suppress_string_rule() {
    let context = TestContext::new();
    let config = r#"
[[hooks.pre-commit]]
type = "suppress-string"
regex = "TODO"
"#;
    context.setup_and_install_old(config);

    // Create a clean file and commit - should succeed
    context.repo.create_file("clean.txt", "No tasks here");
    context.repo.git(&["add", "clean.txt"]);
    let output = context.repo.commit("commit clean file");
    assert!(output.status.success(), "Commit should succeed for clean file");

    // Create a file with suppressed string and try to commit - should fail
    context.repo.create_file("dirty.txt", "I have a TODO here");
    context.repo.git(&["add", "dirty.txt"]);
    let output = context.repo.commit("commit dirty file");
    assert!(!output.status.success(), "Commit should fail for file with suppressed string");
    
    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(stderr.contains("The following files contain suppressed string: dirty.txt"));
}

#[test]
fn test_suppress_string_with_glob_rule() {
    let context = TestContext::new();
    let config = r#"
[[hooks.pre-commit]]
type = "suppress-string"
regex = "DEBUG"
glob = "*.rs"
"#;
    context.setup_and_install_old(config);

    // Create a .txt file with DEBUG - should succeed because glob is *.rs
    context.repo.create_file("debug.txt", "DEBUG logging");
    context.repo.git(&["add", "debug.txt"]);
    let output = context.repo.commit("commit debug txt");
    assert!(output.status.success(), "Commit should succeed for .txt file even with DEBUG");

    // Create a .rs file with DEBUG - should fail
    context.repo.create_file("main.rs", "println!(\"DEBUG\");");
    context.repo.git(&["add", "main.rs"]);
    let output = context.repo.commit("commit debug rs");
    assert!(!output.status.success(), "Commit should fail for .rs file with DEBUG");
    
    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(stderr.contains("The following files contain suppressed string: main.rs"));
}

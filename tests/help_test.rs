use crate::common::{FishermanBinary, GitTestRepo};

mod common;

#[test]
fn no_parameters() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&[], repo.path());

    assert!(
        !output.status.success(),
        "No parameters are not allowed, exit code should be non-zero"
    );

    assert_help_output(&output.stderr);
}

#[test]
fn help() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&["help"], repo.path());

    assert!(
        output.status.success(),
        "Help command should be executed successfully"
    );

    assert_help_output(&output.stdout);
}

#[test]
fn double_slash_help() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&["--help"], repo.path());

    assert!(
        output.status.success(),
        "Help command should be executed successfully"
    );

    assert_help_output(&output.stdout);
}

#[test]
fn slash_h() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&["-h"], repo.path());

    assert!(
        output.status.success(),
        "Help command should be executed successfully"
    );

    assert_help_output(&output.stdout);
}

fn assert_help_output(output: &Vec<u8>) {
    let output_string = String::from_utf8_lossy(output).to_string();

    assert!(
        output_string.contains(&format!("Version: {}", env!("CARGO_PKG_VERSION"))),
        "Output must contain application version"
    );

    assert!(
        output_string.contains("install"),
        "Should have info about install command"
    );
    assert!(
        output_string.contains("handle"),
        "Should have info about handle command"
    );
    assert!(
        output_string.contains("explain"),
        "Should have info about explain command"
    );
    assert!(
        output_string.contains("help"),
        "Should have info about help command"
    );
}

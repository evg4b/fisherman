use std::env;

use crate::common::{FishermanBinary, GitTestRepo};

mod common;

#[test]
#[cfg(feature = "integration-tests")]
fn double_slash_version() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&["--version"], repo.path());

    assert!(output.status.success(), "Should complete successfully");
    let output_string = String::from_utf8_lossy(&output.stdout).to_string();

    assert!(
        output_string.contains(env!("CARGO_PKG_VERSION")),
        "Output must contain application version"
    );
}

#[test]
#[cfg(feature = "integration-tests")]
fn slash_v() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let output = fisherman.run(&["-V"], repo.path());

    assert!(output.status.success(), "Should complete successfully");
    let output_string = String::from_utf8_lossy(&output.stdout).to_string();

    assert!(
        output_string.contains(env!("CARGO_PKG_VERSION")),
        "Output must contain application version"
    );
}

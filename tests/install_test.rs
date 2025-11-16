mod common;

use tempdir::TempDir;

use crate::common::{FishermanBinary, GitTestRepo, hooks::Hook};


#[test]
fn install_in_empty_dir() {
    let temp_dir = TempDir::new("fisherman_test").expect("Failed to create temp directory");
    let fisherman = FishermanBinary::build();
    
    let output = fisherman.install(temp_dir.path(), false);
    
    assert!(!output.status.success());
    assert!(String::from_utf8_lossy(&output.stderr).contains("Error: could not find repository at "));
    assert!(String::from_utf8_lossy(&output.stderr).contains(temp_dir.path().to_str().unwrap()));
}

#[test]
fn install_in_empty_repo() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let ouput = fisherman.install(repo.path(), false);

    assert!(ouput.status.success());
    for hook in Hook::iter() {
        assert!(String::from_utf8_lossy(&ouput.stdout)
            .contains(hook.as_str()))
    }
}
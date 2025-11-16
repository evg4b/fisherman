mod common;

use tempdir::TempDir;

use crate::common::{FishermanBinary, GitTestRepo, hooks::Hook};

#[test]
fn install_in_empty_dir() {
    let temp_dir = TempDir::new("fisherman_test").expect("Failed to create temp directory");
    let fisherman = FishermanBinary::build();
    
    let output = fisherman.install(temp_dir.path(), false);
    
    assert!(!output.status.success());
    assert!(
        String::from_utf8_lossy(&output.stderr)
            .contains("Error: could not find repository at ")
    );
    assert!(
        String::from_utf8_lossy(&output.stderr)
            .contains(temp_dir.path().to_str().unwrap())
    );
}

#[test]
fn install_in_empty_repo() {
    let repo = GitTestRepo::new();
    let fisherman = FishermanBinary::build();

    let ouput = fisherman.install(repo.path(), false);

    assert!(ouput.status.success());
    for hook in Hook::iter() {
        assert!(repo.hook_exists(hook.as_str()));
        assert!(String::from_utf8_lossy(&ouput.stdout)
            .contains(hook.as_str()));
    }
}

#[test]
fn install_using_local_conig() {
    let repo = GitTestRepo::new();

    let fisherman = FishermanBinary::build();

    repo.create_config(r#"[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"

[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(feature|bugfix)/[a-zA-Z0-9-_]+$"
"#);

    let ouput = fisherman.install(repo.path(), false);

    assert!(ouput.status.success());
    for hook in Hook::iter() {
        if *hook == Hook::PrePush || *hook == Hook::CommitMsg {
            assert!(repo.hook_exists(hook.as_str()));
            assert!(String::from_utf8_lossy(&ouput.stdout)
                .contains(hook.as_str()));
        } else {
            assert!(!repo.hook_exists(hook.as_str()));
            assert!(!String::from_utf8_lossy(&ouput.stdout)
                .contains(hook.as_str()));
        }
    }
    
}
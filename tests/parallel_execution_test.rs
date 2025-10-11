mod common;

use common::{FishermanBinary, GitTestRepo};
use std::time::Instant;

#[test]
fn parallel_multiple_write_files() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "file1.txt"
content = "content 1"

[[hooks.pre-commit]]
type = "write-file"
path = "file2.txt"
content = "content 2"

[[hooks.pre-commit]]
type = "write-file"
path = "file3.txt"
content = "content 3"

[[hooks.pre-commit]]
type = "write-file"
path = "file4.txt"
content = "content 4"

[[hooks.pre-commit]]
type = "write-file"
path = "file5.txt"
content = "content 5"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "All write-file rules should execute successfully"
    );

    assert!(repo.file_exists("file1.txt"));
    assert!(repo.file_exists("file2.txt"));
    assert!(repo.file_exists("file3.txt"));
    assert!(repo.file_exists("file4.txt"));
    assert!(repo.file_exists("file5.txt"));
}

#[test]
fn parallel_multiple_exec_rules() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "1"]

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "2"]

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "3"]
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["1"]

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["2"]

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["3"]
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "All exec rules should execute successfully"
    );
}

#[test]
fn parallel_multiple_shell_scripts() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "echo script1"

[[hooks.pre-commit]]
type = "shell"
script = "echo script2"

[[hooks.pre-commit]]
type = "shell"
script = "echo script3"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "echo 'script1'"

[[hooks.pre-commit]]
type = "shell"
script = "echo 'script2'"

[[hooks.pre-commit]]
type = "shell"
script = "echo 'script3'"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "All shell scripts should execute successfully"
    );
}

#[test]
fn parallel_mixed_async_rules() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output1.txt"
content = "write-file"

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "exec"]

[[hooks.pre-commit]]
type = "shell"
script = "echo shell"

[[hooks.pre-commit]]
type = "write-file"
path = "output2.txt"
content = "another write"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "output1.txt"
content = "write-file"

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["exec"]

[[hooks.pre-commit]]
type = "shell"
script = "echo 'shell'"

[[hooks.pre-commit]]
type = "write-file"
path = "output2.txt"
content = "another write"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "All async rules should execute successfully"
    );
    assert!(repo.file_exists("output1.txt"));
    assert!(repo.file_exists("output2.txt"));
}

#[test]
fn parallel_one_fails_stops_all() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "success"]

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "exit", "1"]

[[hooks.pre-commit]]
type = "write-file"
path = "should-not-exist.txt"
content = "should not be created"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["success"]

[[hooks.pre-commit]]
type = "exec"
command = "false"

[[hooks.pre-commit]]
type = "write-file"
path = "should-not-exist.txt"
content = "should not be created"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        !handle_output.status.success(),
        "Hook should fail when one async rule fails"
    );
}

#[test]
fn sync_rules_execute_before_async() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "async"]

[[hooks.pre-commit]]
type = "write-file"
path = "async.txt"
content = "async rule"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["async"]

[[hooks.pre-commit]]
type = "write-file"
path = "async.txt"
content = "async rule"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/test");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Sync rules should execute before async rules"
    );
    assert!(repo.file_exists("async.txt"));
}

#[test]
fn sync_rule_fails_hook_fails() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/test");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        !handle_output.status.success(),
        "Hook should fail when sync rule fails"
    );
}

#[test]
#[cfg(not(windows))]
fn parallel_performance_benefit() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "sleep 0.1"

[[hooks.pre-commit]]
type = "shell"
script = "sleep 0.1"

[[hooks.pre-commit]]
type = "shell"
script = "sleep 0.1"

[[hooks.pre-commit]]
type = "shell"
script = "sleep 0.1"

[[hooks.pre-commit]]
type = "shell"
script = "sleep 0.1"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let start = Instant::now();
    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    let duration = start.elapsed();

    assert!(handle_output.status.success());
    assert!(
        duration.as_millis() < 1000,
        "Parallel execution should take less than 1000ms (sequential would be 500ms+), took {}ms",
        duration.as_millis()
    );
}

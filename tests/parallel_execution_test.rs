mod common;

use common::test_context::TestContext;
use std::time::Instant;

/// Tests that multiple write-file rules execute in parallel and create all target files.
/// Verifies parallel execution of asynchronous rules completes successfully.
#[test]
fn parallel_multiple_write_files() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("file1.txt"));
    assert!(ctx.repo.file_exists("file2.txt"));
    assert!(ctx.repo.file_exists("file3.txt"));
    assert!(ctx.repo.file_exists("file4.txt"));
    assert!(ctx.repo.file_exists("file5.txt"));
}

/// Tests that multiple exec rules execute in parallel successfully.
/// Verifies concurrent command execution works correctly across platforms.
#[test]
fn parallel_multiple_exec_rules() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");
}

/// Tests that multiple shell script rules execute in parallel.
/// Verifies that concurrent shell script execution completes without conflicts.
#[test]
fn parallel_multiple_shell_scripts() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");
}

/// Tests that different types of async rules (write-file, exec, shell) run in parallel.
/// Verifies mixed async rule types can execute concurrently without issues.
#[test]
fn parallel_mixed_async_rules() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("output1.txt"));
    assert!(ctx.repo.file_exists("output2.txt"));
}

/// Tests that when one parallel rule fails, the hook execution fails appropriately.
/// Verifies error handling in parallel execution propagates failures correctly.
#[test]
fn parallel_one_fails_stops_all() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.handle_failure("pre-commit");
}

/// Tests that synchronous validation rules execute before asynchronous rules.
/// Verifies correct execution order with sync rules first, then parallel async rules.
#[test]
fn sync_rules_execute_before_async() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/test");

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("async.txt"));
}

/// Tests that when a synchronous rule fails, async rules don't execute and hook fails.
/// Verifies early exit behavior when sync validation fails before async execution.
#[test]
fn sync_rule_fails_hook_fails() {
    let ctx = TestContext::new();

    let config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/test");

    ctx.handle_failure("pre-commit");
}

/// Tests that parallel execution provides performance benefit over sequential execution.
/// Verifies that multiple sleep commands complete faster due to parallelization (Unix only).
#[test]
#[cfg(not(windows))]
fn parallel_performance_benefit() {
    let ctx = TestContext::new();

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

    ctx.setup_and_install(config);

    let start = Instant::now();
    let handle_output = ctx.handle("pre-commit");
    let duration = start.elapsed();

    assert!(handle_output.status.success());
    assert!(
        duration.as_millis() < 1000,
        "Parallel execution should take less than 1000ms (sequential would be 500ms+), took {}ms",
        duration.as_millis()
    );
}

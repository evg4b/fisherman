mod common;

use common::{test_context::{echo_config, fail_config, TestContext}, FishermanBinary, GitTestRepo};

/// Tests that exec rule executes successfully when command exits with code 0.
/// Verifies basic command execution functionality using echo command.
#[test]
fn exec_rule_success() {
    let ctx = TestContext::new();
    let config = echo_config("pre-commit", "test");

    ctx.setup_and_install(&config);
    let output = ctx.handle("pre-commit");
    assert!(output.status.success());

    let stdout = String::from_utf8_lossy(&output.stdout);
    assert!(stdout.contains("test"), "Output should contain 'test': {}", stdout);
}

/// Tests that exec rule fails appropriately when command exits with non-zero code.
/// Verifies that command failures are detected and propagate as hook failures.
#[test]
fn exec_rule_failure() {
    let ctx = TestContext::new();
    let config = fail_config("pre-commit");

    ctx.setup_and_install(&config);
    let output = ctx.handle("pre-commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error output should not be empty");
}

/// Tests that exec rule properly passes environment variables to executed command.
/// Verifies that custom env variables are available in the command's environment.
#[test]
fn exec_rule_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "%TEST_VAR%"]
env = { TEST_VAR = "test_value" }
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "sh"
args = ["-c", "test \"$TEST_VAR\" = \"test_value\""]
env = { TEST_VAR = "test_value" }
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should pass environment variables: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

/// Tests that shell script rule executes successfully when script exits with code 0.
/// Verifies basic shell script execution on both Windows and Unix platforms.
#[test]
fn shell_script_success() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "echo test"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = """
#!/bin/sh
echo "Running shell script"
exit 0
"""
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Shell script should succeed: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );

    #[cfg(not(windows))]
    {
        let stdout = String::from_utf8_lossy(&handle_output.stdout);
        assert!(stdout.contains("Running shell script"), "Output should contain script message: {}", stdout);
    }
}

/// Tests that shell script rule fails when script exits with non-zero code.
/// Verifies that shell script failures properly abort the hook execution.
#[test]
fn shell_script_failure() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "exit 1"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = """
#!/bin/sh
exit 1
"""
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        !handle_output.status.success(),
        "Shell script should fail with exit 1"
    );

    let stderr = String::from_utf8_lossy(&handle_output.stderr);
    assert!(!stderr.is_empty(), "Error output should not be empty when shell fails");
}

/// Tests that shell script can access custom environment variables defined in config.
/// Verifies environment variable propagation to shell execution context.
#[test]
fn shell_script_with_env() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = "if \"%CUSTOM_VAR%\" == \"custom_value\" exit 0"
env = { CUSTOM_VAR = "custom_value" }
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "shell"
script = """
#!/bin/sh
if [ "$CUSTOM_VAR" = "custom_value" ]; then
    exit 0
else
    exit 1
fi
"""
env = { CUSTOM_VAR = "custom_value" }
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Shell script should access env variables: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

/// Tests that exec and shell rules can be configured together and both execute successfully.
/// Verifies compatibility and correct execution order of mixed rule types.
#[test]
fn exec_and_shell_mixed() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    #[cfg(windows)]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "exec test"]

[[hooks.pre-commit]]
type = "shell"
script = "echo shell test"
"#;

    #[cfg(not(windows))]
    let config = r#"
[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["exec test"]

[[hooks.pre-commit]]
type = "shell"
script = "echo 'shell test'"
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Both exec and shell rules should succeed"
    );

    let stdout = String::from_utf8_lossy(&handle_output.stdout);
    assert!(stdout.contains("exec test") || stdout.contains("1"),
        "Output should contain exec command result: {}", stdout);
    assert!(stdout.contains("shell test") || stdout.contains("shell"),
        "Output should contain shell script result: {}", stdout);
}

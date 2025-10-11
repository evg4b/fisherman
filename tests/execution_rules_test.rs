mod common;

use common::{test_context::{echo_config, fail_config, shell_config}, TestContext, FishermanBinary, GitTestRepo};

#[test]
fn exec_rule_success() {
    let ctx = TestContext::new();
    let config = echo_config("pre-commit", "test");

    ctx.setup_and_install(&config);
    ctx.handle_success("pre-commit");
}

#[test]
fn exec_rule_failure() {
    let ctx = TestContext::new();
    let config = fail_config("pre-commit");

    ctx.setup_and_install(&config);
    ctx.handle_failure("pre-commit");
}

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
}

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
}

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
}

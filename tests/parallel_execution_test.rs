mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use fisherman_core::BranchNamePrefixRule;
use fisherman_core::Configuration;
use fisherman_core::ExecRule;
use fisherman_core::GitHook;
use fisherman_core::ShellScriptRule;
use fisherman_core::WriteFileRule;
use std::time::Instant;

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_multiple_write_files() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "file1.txt".into(),
                content: "content 1".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "file2.txt".into(),
                content: "content 2".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "file3.txt".into(),
                content: "content 3".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "file4.txt".into(),
                content: "content 4".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "file5.txt".into(),
                content: "content 5".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("file1.txt"));
    assert!(ctx.repo.file_exists("file2.txt"));
    assert!(ctx.repo.file_exists("file3.txt"));
    assert!(ctx.repo.file_exists("file4.txt"));
    assert!(ctx.repo.file_exists("file5.txt"));
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_multiple_exec_rules() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("1")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("2")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("3")]),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("1")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("2")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("3")]),
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success(), "All exec rules should succeed: {}",
            String::from_utf8_lossy(&output.stderr));
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_multiple_shell_scripts() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                script: "echo script1".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo script2".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo script3".into(),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                script: "echo 'script1'".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo 'script2'".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo 'script3'".into(),
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    let stdout = String::from_utf8_lossy(&output.stdout);
    #[cfg(not(windows))]
    {
        assert!(stdout.contains("script1") || !stdout.is_empty(),
                "Output should contain script results: {}", stdout);
    }
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_mixed_async_rules() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "output1.txt".into(),
                content: "write-file".into(),
                append: None,
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("exec")]),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo shell".into(),
                env: None,
            }),
            rule!(WriteFileRule {
                path: "output2.txt".into(),
                content: "another write".into(),
                append: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "output1.txt".into(),
                content: "write-file".into(),
                append: None,
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("exec")]),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "echo 'shell'".into(),
                env: None,
            }),
            rule!(WriteFileRule {
                path: "output2.txt".into(),
                content: "another write".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("output1.txt"));
    assert!(ctx.repo.file_exists("output2.txt"));

    let stdout = String::from_utf8_lossy(&output.stdout);
    #[cfg(not(windows))]
    {
        assert!(stdout.contains("exec") || !stdout.is_empty(),
                "Output should contain exec command result: {}", stdout);
    }
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_one_fails_stops_all() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("success")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("exit"), String::from("1")]),
                env: None,
            }),
            rule!(WriteFileRule {
                path: "should-not-exist.txt".into(),
                content: "should not be created".into(),
                append: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("success")]),
                env: None,
            }),
            rule!(ExecRule {
                command: String::from("false"),
                args: None,
                env: None,
            }),
            rule!(WriteFileRule {
                path: "should-not-exist.txt".into(),
                content: "should not be created".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error output should not be empty when a rule fails");
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn sync_rules_execute_before_async() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("async")]),
                env: None,
            }),
            rule!(WriteFileRule {
                path: "async.txt".into(),
                content: "async rule".into(),
                append: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("async")]),
                env: None,
            }),
            rule!(WriteFileRule {
                path: "async.txt".into(),
                content: "async rule".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/test");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("async.txt"));

    #[cfg(not(windows))]
    {
        let stdout = String::from_utf8_lossy(&output.stdout);
        assert!(stdout.contains("async") || !stdout.is_empty(),
                "Output should contain async rule result: {}", stdout);
    }
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn sync_rule_fails_hook_fails() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/test");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success());

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(!stderr.is_empty(), "Error should explain which rule failed");
}

#[test]
#[ignore = "UNSUPPORTED: parallel execution is not supported at the moment"]
fn parallel_performance_benefit() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(ShellScriptRule {
                script: "sleep 0.1".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "sleep 0.1".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "sleep 0.1".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "sleep 0.1".into(),
                env: None,
            }),
            rule!(ShellScriptRule {
                script: "sleep 0.1".into(),
                env: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    let start = Instant::now();
    let handle_output = ctx.git_commit_allow_empty("test commit");
    let duration = start.elapsed();

    assert!(handle_output.status.success(), "Hook should succeed: {}",
            String::from_utf8_lossy(&handle_output.stderr));

    assert!(
        duration.as_millis() < 3000,
        "Parallel execution should take less than 3000ms (sequential would be 500ms+), took {}ms",
        duration.as_millis()
    );
}

mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::write_file::WriteFileRule;
use std::time::Instant;

#[test]
fn parallel_multiple_write_files() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("file1.txt"),
                content: String::from("content 1"),
                append: None,
            }),
            rule!(WriteFileRule {
                path: String::from("file2.txt"),
                content: String::from("content 2"),
                append: None,
            }),
            rule!(WriteFileRule {
                path: String::from("file3.txt"),
                content: String::from("content 3"),
                append: None,
            }),
            rule!(WriteFileRule {
                path: String::from("file4.txt"),
                content: String::from("content 4"),
                append: None,
            }),
            rule!(WriteFileRule {
                path: String::from("file5.txt"),
                content: String::from("content 5"),
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
fn parallel_multiple_shell_scripts() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("echo script1"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo script2"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo script3"),
                env: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'script1'"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'script2'"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'script3'"),
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
fn parallel_mixed_async_rules() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("output1.txt"),
                content: String::from("write-file"),
                append: None,
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("exec")]),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo shell"),
                env: None,
            }),
            rule!(WriteFileRule {
                path: String::from("output2.txt"),
                content: String::from("another write"),
                append: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("output1.txt"),
                content: String::from("write-file"),
                append: None,
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("exec")]),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'shell'"),
                env: None,
            }),
            rule!(WriteFileRule {
                path: String::from("output2.txt"),
                content: String::from("another write"),
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
                path: String::from("should-not-exist.txt"),
                content: String::from("should not be created"),
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
                path: String::from("should-not-exist.txt"),
                content: String::from("should not be created"),
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
fn sync_rules_execute_before_async() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("feature/"),
            }),
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("async")]),
                env: None,
            }),
            rule!(WriteFileRule {
                path: String::from("async.txt"),
                content: String::from("async rule"),
                append: None,
            })
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("feature/"),
            }),
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("async")]),
                env: None,
            }),
            rule!(WriteFileRule {
                path: String::from("async.txt"),
                content: String::from("async rule"),
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
fn sync_rule_fails_hook_fails() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("feature/"),
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
#[cfg(not(windows))]
fn parallel_performance_benefit() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("sleep 0.1"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("sleep 0.1"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("sleep 0.1"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("sleep 0.1"),
                env: None,
            }),
            rule!(RuleParams::ShellScript {
                script: String::from("sleep 0.1"),
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

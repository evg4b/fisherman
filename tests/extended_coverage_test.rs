mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use fisherman_core::{t, CommitMessageRegexRule, BranchNameRegexRule, BranchNamePrefixRule, BranchNameSuffixRule, Configuration, Expression, ExecRule, GitHook, ShellScriptRule, WriteFileRule};
use std::collections::HashMap;

#[test]
fn unicode_in_commit_message() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: "^.+$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("feat: Add 日本語 support with émojis 🎉");
}

#[test]
fn unicode_in_branch_name() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameRegexRule {
                expression: "^.+$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/support-日本語");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn unicode_in_template_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "branch-name.txt".into(),
                content: t!("Branch: {{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Name>.+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/日本語-support");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("branch-name.txt");
    assert!(content.contains("日本語"));
}

#[test]
fn prepare_commit_msg_hook_execution() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PrepareCommitMsg => [
            rule!(WriteFileRule {
                path: "prepare-executed.txt".into(),
                content: "prepare-commit-msg ran".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("prepare-executed.txt"));
}

#[test]
fn prepare_commit_msg_with_template_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PrepareCommitMsg => [
            rule!(WriteFileRule {
                path: "commit-template.txt".into(),
                content: t!("{{Type}}: [{{Ticket}}] "),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-456");

    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("commit-template.txt"));
    let content = ctx.repo.read_file("commit-template.txt");
    assert_eq!(content, "feature: [PROJ-456] ");
}

#[test]
fn conditional_with_undefined_variable_fails() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "output.txt".into(),
                content: "test".into(),
                append: None,
            }, when = Expression::new("UndefinedVar == \"value\""))
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_failure("test commit");
}

#[test]
fn conditional_with_is_def_var_true() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "output.txt".into(),
                content: "Feature is defined".into(),
                append: None,
            }, when = Expression::new("is_def_var(\"Feature\")"))
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth");
    ctx.git_commit_allow_empty_success("test commit");

    assert_eq!(ctx.repo.read_file("output.txt"), "Feature is defined");
}

#[test]
fn conditional_with_is_def_var_false() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "output.txt".into(),
                content: "Feature is defined".into(),
                append: None,
            }, when = Expression::new("is_def_var(\"Feature\")")),
            rule!(WriteFileRule {
                path: "fallback.txt".into(),
                content: "Feature not defined".into(),
                append: None,
            }, when = Expression::new("!is_def_var(\"Feature\")"))
        ],
        extract = vec![
            String::from("branch?:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/test");
    ctx.git_commit_allow_empty_success("test commit");

    assert!(!ctx.repo.file_exists("output.txt"));
    assert!(ctx.repo.file_exists("fallback.txt"));
}

#[test]
fn shell_script_with_multiple_env_vars() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = {
        let mut env = HashMap::new();
        env.insert(String::from("VAR1"), String::from("value1"));
        env.insert(String::from("VAR2"), String::from("value2"));

        config!(
            GitHook::PreCommit => [
                rule!(ShellScriptRule {
                    script: "if \"%VAR1%\" == \"value1\" if \"%VAR2%\" == \"value2\" exit 0".into(),
                    env: Some(env),
                })
            ]
        )
    };

    #[cfg(not(windows))]
    let config = {
        let mut env = HashMap::new();
        env.insert(String::from("VAR1"), String::from("value1"));
        env.insert(String::from("VAR2"), String::from("value2"));

        config!(
            GitHook::PreCommit => [
                rule!(ShellScriptRule {
                    script: "#!/bin/sh\nif [ \"$VAR1\" = \"value1\" ] && [ \"$VAR2\" = \"value2\" ]; then\n    exit 0\nelse\n    exit 1\nfi".into(),
                    env: Some(env),
                })
            ]
        )
    };

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn exec_with_templated_env_vars() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = {
        let mut env = HashMap::new();
        env.insert(String::from("FEATURE_NAME"), String::from("{{Feature}}"));

        config!(
            GitHook::PreCommit => [
                rule!(ExecRule {
                    command: "cmd".into(),
                    args: Some(vec!["/C".to_string(), "echo".to_string(), "%FEATURE_NAME%".to_string()]),
                    env: Some(env),
                })
            ],
            extract = vec![
                String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
            ]
        )
    };

    #[cfg(not(windows))]
    let config = {
        let mut env = HashMap::new();
        env.insert(String::from("FEATURE_NAME"), String::from("{{Feature}}"));

        config!(
            GitHook::PreCommit => [
                rule!(ExecRule {
                    command: "sh".into(),
                    args: Some(vec!["-c".to_string(), "test \"$FEATURE_NAME\" = \"payment\"".to_string()]),
                    env: Some(env),
                })
            ],
            extract = vec![
                String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
            ]
        )
    };

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/payment");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn empty_commit_message() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: "^.+$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_failure("");
}

#[test]
fn very_long_commit_message() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: "^.+$".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let long_message = "a".repeat(10000);
    ctx.git_commit_allow_empty_success(&long_message);
}

#[test]
fn whitespace_only_commit_message_rejected_by_git() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: ".*".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    let output = ctx.git_commit_allow_empty("   \n   \t   ");
    assert!(!output.status.success(), "Git should reject whitespace-only commit messages");
}

#[test]
fn write_file_with_special_characters_in_content() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "special.txt".into(),
                content: "Line with $VAR and `backticks` and \"quotes\" and 'apostrophes'".into(),
                append: None,
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("special.txt");
    assert!(content.contains("$VAR"));
    assert!(content.contains("`backticks`"));
    assert!(content.contains("\"quotes\""));
}

#[test]
fn branch_name_with_slashes() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "hierarchy.txt".into(),
                content: t!("{{Category}}/{{Subcategory}}/{{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Category>[^/]+)/(?P<Subcategory>[^/]+)/(?P<Name>.+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/ui/button-component");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("hierarchy.txt");
    assert_eq!(content, "feature/ui/button-component");
}

#[test]
fn multiple_rules_with_mixed_success_sync() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            }),
            rule!(BranchNameRegexRule {
                expression: "^feature/[a-z-]+$".into(),
            }),
            rule!(BranchNameSuffixRule {
                suffix: "-ready".into(),
            })
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/new-feature-ready");
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn regex_with_escaped_characters() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "ticket.txt".into(),
                content: t!("{{Ticket}} - {{Priority}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Ticket>[A-Z]+-\\d+)-(?P<Priority>high|low)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-123-high");
    ctx.git_commit_allow_empty_success("test commit");

    let content = ctx.repo.read_file("ticket.txt");
    assert_eq!(content, "PROJ-123 - high");
}

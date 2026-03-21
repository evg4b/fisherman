mod common;

use common::test_context::TestContext;
use common::ConfigFormat;
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::RuleParams;

#[test]
fn template_branch_variable_in_message_prefix() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessagePrefix {
                prefix: String::from("{{Type}}: [{{Ticket}}] "),
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.git_commit_allow_empty_success("feature: [PROJ-123] add new feature");
}

#[test]
fn template_branch_variable_in_write_file() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("branch-info.txt"),
                content: String::from("Current feature: {{Feature}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth-system");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("branch-info.txt"));
    assert_eq!(
        ctx.repo.read_file("branch-info.txt"),
        "Current feature: auth-system"
    );
}

#[test]
fn template_repo_path_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("repo-info.txt"),
                content: String::from("Repository: {{RepoName}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("repo_path:.*/(?P<RepoName>[^/]+)$"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("repo-info.txt"));

    let content = ctx.repo.read_file("repo-info.txt");
    assert!(content.starts_with("Repository: "));
}

#[test]
fn template_multiple_variables() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("info.txt"),
                content: String::from("Type: {{Type}}, Ticket: {{Ticket}}, Repo: {{RepoName}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"),
            String::from("repo_path:.*/(?P<RepoName>[^/]+)$"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/ABC-456");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("info.txt"));

    let content = ctx.repo.read_file("info.txt");
    assert!(content.contains("Type: feature"));
    assert!(content.contains("Ticket: ABC-456"));
    assert!(content.contains("Repo: "));
}

#[test]
fn template_in_exec_command() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("{{Feature}}")]),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("{{Feature}}")]),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/payment");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_optional_variable_present() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("output.txt"),
                content: String::from("Feature: {{Feature}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch?:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth");

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "Feature: auth");
}

#[test]
fn template_optional_variable_missing() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNameRegex {
                regex: String::from("^.+$"),
            })
        ],
        extract = vec![
            String::from("branch?:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/issue");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_in_file_path() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("{{Feature}}-status.txt"),
                content: String::from("Feature status file"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/payment");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("payment-status.txt"));
    assert_eq!(
        ctx.repo.read_file("payment-status.txt"),
        "Feature status file"
    );
}

#[test]
fn template_in_file_path_multiple_vars() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("{{Type}}-{{Name}}.log"),
                content: String::from("Log for {{Type}}/{{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Name>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("feature-auth.log"));
    assert_eq!(
        ctx.repo.read_file("feature-auth.log"),
        "Log for feature/auth"
    );
}

#[test]
fn template_in_message_suffix() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::CommitMsg => [
            rule!(RuleParams::CommitMessageSuffix {
                suffix: String::from(" [{{Ticket}}]"),
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Ticket>[A-Z]+-\\d+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/PROJ-789");

    ctx.git_commit_allow_empty_success("Add new feature [PROJ-789]");
    ctx.git_commit_allow_empty_failure("Add new feature");
}

#[test]
fn template_in_branch_name_prefix() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PrePush => [
            rule!(RuleParams::BranchNamePrefix {
                prefix: String::from("{{Prefix}}/"),
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Prefix>feature|bugfix)/"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth-system");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_in_branch_name_suffix() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PrePush => [
            rule!(RuleParams::BranchNameSuffix {
                suffix: String::from("-{{Type}}"),
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth-feature");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_in_shell_command() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("echo Working on {{Feature}} > feature.txt"),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ShellScript {
                script: String::from("echo 'Working on {{Feature}}' > feature.txt"),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/dashboard");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("feature.txt"));

    let content = ctx.repo.read_file("feature.txt");
    assert!(content.contains("dashboard"));
}

#[test]
fn multiple_templates_in_single_field() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("status.txt"),
                content: String::from("Type: {{Type}}, Ticket: {{Ticket}}, Name: {{Name}}, Full: {{Type}}/{{Ticket}}-{{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)-(?P<Name>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/ABC-123-login");

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(
        ctx.repo.read_file("status.txt"),
        "Type: feature, Ticket: ABC-123, Name: login, Full: feature/ABC-123-login"
    );
}

#[test]
fn template_rendering_failure_missing_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("output.txt"),
                content: String::from("Feature: {{Feature}}, Missing: {{UndefinedVar}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth");

    ctx.git_commit_allow_empty_failure("test commit");
}

#[test]
fn template_optional_repo_path_no_match() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::BranchNameRegex {
                regex: String::from(".*"),
            })
        ],
        extract = vec![
            String::from("repo_path?:^/nonexistent/(?P<Project>[^/]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_with_special_characters() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("feature.txt"),
                content: String::from("Feature: {{Feature}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z0-9_-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth_v2-beta");

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(
        ctx.repo.read_file("feature.txt"),
        "Feature: auth_v2-beta"
    );
}

#[test]
fn template_in_conditional_with_defined_var() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("conditional.txt"),
                    content: String::from("Type: {{Type}}"),
                    append: None,
                },
                when = String::from("is_def_var(\"Type\")")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Name>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/auth");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("conditional.txt"));
}

#[test]
fn template_conditional_skipped_undefined_var() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("optional.txt"),
                    content: String::from("Feature: {{Feature}}"),
                    append: None,
                },
                when = String::from("is_def_var(\"Feature\")")
            )
        ],
        extract = vec![
            String::from("branch?:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("bugfix/issue");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(!ctx.repo.file_exists("optional.txt"));
}

#[test]
fn template_complex_extraction_pattern() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("{{Category}}-{{Project}}-{{Issue}}.txt"),
                content: String::from("Category: {{Category}}\nProject: {{Project}}\nIssue: {{Issue}}\nDescription: {{Description}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Category>feature|bugfix|hotfix)/(?P<Project>[A-Z]+)-(?P<Issue>\\d+)-(?P<Description>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/MYAPP-456-user-auth");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("feature-MYAPP-456.txt"));

    let content = ctx.repo.read_file("feature-MYAPP-456.txt");
    assert!(content.contains("Category: feature"));
    assert!(content.contains("Project: MYAPP"));
    assert!(content.contains("Issue: 456"));
    assert!(content.contains("Description: user-auth"));
}

#[test]
fn template_in_multiple_exec_args() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("cmd"),
                args: Some(vec![String::from("/C"), String::from("echo"), String::from("{{Type}}"), String::from("{{Name}}")]),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Name>[a-z-]+)"),
        ]
    );

    #[cfg(not(windows))]
    let config = config!(
        GitHook::PreCommit => [
            rule!(ExecRule {
                command: String::from("echo"),
                args: Some(vec![String::from("{{Type}}"), String::from("{{Name}}")]),
                env: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Name>[a-z-]+)"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/dashboard");

    ctx.git_commit_allow_empty_success("test commit");
}

#[test]
fn template_combined_repo_and_branch_variables() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: String::from("combined-info.log"),
                content: String::from("Repo: {{RepoName}}, Type: {{Type}}, Name: {{Name}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Name>[a-z-]+)"),
            String::from("repo_path:.*/(?P<RepoName>[^/]+)$"),
        ]
    );

    ctx.setup_and_install(&config, ConfigFormat::Toml);
    ctx.repo.create_branch("feature/api");

    ctx.git_commit_allow_empty_success("test commit");

    assert!(ctx.repo.file_exists("combined-info.log"));
    let content = ctx.repo.read_file("combined-info.log");
    assert!(content.contains("Type: feature"));
    assert!(content.contains("Name: api"));
    assert!(content.contains("Repo: "));
}


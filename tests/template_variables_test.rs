mod common;

use common::configuration::serialize_configuration;
use common::test_context::TestContext;
use common::ConfigFormat;
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::RuleParams;

/// Tests that template variables extracted from branch name work in message-prefix rule.
/// Verifies variable extraction and template rendering for commit message validation.
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

/// Tests that template variables from branch name are correctly substituted in write-file content.
/// Verifies variable extraction and file content rendering with extracted values.
#[test]
fn template_branch_variable_in_write_file() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that repository path can be extracted and used in template variables.
/// Verifies repo_path extraction pattern works and renders in file content.
#[test]
fn template_repo_path_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that multiple variables from different sources can be extracted and used together.
/// Verifies simultaneous branch and repo_path variable extraction and template rendering.
#[test]
fn template_multiple_variables() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that template variables work in exec command arguments.
/// Verifies variable substitution in command-line arguments for exec rules.
#[test]
fn template_in_exec_command() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
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
            rule!(RuleParams::ExecRule {
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

/// Tests that optional variables (branch?) are extracted when pattern matches.
/// Verifies optional extraction syntax works correctly when branch name matches pattern.
#[test]
fn template_optional_variable_present() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that optional variables (branch?) don't cause failure when pattern doesn't match.
/// Verifies that optional extraction allows hook to proceed even when variable isn't extracted.
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

/// Tests that template variables work in file paths for write-file rules.
/// Verifies dynamic file naming based on extracted variables.
#[test]
fn template_in_file_path() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that template variables work in flat file paths with multiple variables.
/// Verifies template interpolation in file names with multiple substitutions.
#[test]
fn template_in_file_path_multiple_vars() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that template variables work in message-suffix rules.
/// Verifies variable substitution in commit message suffix validation.
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

/// Tests that template variables work in branch-name-prefix rules.
/// Verifies variable substitution in branch name validation with prefix.
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

/// Tests that template variables work in branch-name-suffix rules.
/// Verifies variable substitution in branch name validation with suffix.
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

/// Tests that template variables work in shell commands.
/// Verifies variable substitution in shell script execution.
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

/// Tests that multiple templates in the same field are all interpolated.
/// Verifies that multiple {{variable}} occurrences work correctly.
#[test]
fn multiple_templates_in_single_field() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that missing required variables cause hook execution to fail.
/// Verifies that undefined variables in templates are properly detected and reported.
#[test]
fn template_rendering_failure_missing_variable() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that optional repo_path extraction doesn't fail when pattern doesn't match.
/// Verifies optional repo_path? syntax allows execution to continue.
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

/// Tests that template variables with special characters are handled correctly.
/// Verifies that hyphens, underscores, and numbers in extracted values work.
#[test]
fn template_with_special_characters() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that templates work in conditional expressions with defined variables.
/// Verifies template interpolation in when conditions combined with is_def_var.
#[test]
fn template_in_conditional_with_defined_var() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                RuleParams::WriteFile {
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

/// Tests that rules are skipped when conditional variables are not defined.
/// Verifies that when conditions properly evaluate variable existence.
#[test]
fn template_conditional_skipped_undefined_var() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                RuleParams::WriteFile {
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

/// Tests complex extraction pattern with multiple capture groups.
/// Verifies that complex regex patterns with multiple groups work correctly.
#[test]
fn template_complex_extraction_pattern() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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

/// Tests that templates work correctly in exec command arguments.
/// Verifies variable substitution across multiple command arguments.
#[test]
fn template_in_multiple_exec_args() {
    let ctx = TestContext::new();

    #[cfg(windows)]
    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::ExecRule {
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
            rule!(RuleParams::ExecRule {
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

/// Tests template interpolation with repo_path and branch combined.
/// Verifies that both extraction sources work together in the same rule.
#[test]
fn template_combined_repo_and_branch_variables() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(RuleParams::WriteFile {
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


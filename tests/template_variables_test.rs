mod common;

use common::test_context::TestContext;

/// Tests that template variables extracted from branch name work in message-prefix rule.
/// Verifies variable extraction and template rendering for commit message validation.
#[test]
fn template_branch_variable_in_message_prefix() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.commit-msg]]
type = "message-prefix"
prefix = "{{Type}}: [{{Ticket}}] "
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/PROJ-123");

    ctx.handle_commit_msg_success("feature: [PROJ-123] add new feature");
}

/// Tests that template variables from branch name are correctly substituted in write-file content.
/// Verifies variable extraction and file content rendering with extracted values.
#[test]
fn template_branch_variable_in_write_file() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "branch-info.txt"
content = "Current feature: {{Feature}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/auth-system");

    ctx.handle_success("pre-commit");
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

    let config = r#"
extract = ["repo_path:.*/(?P<RepoName>[^/]+)$"]

[[hooks.pre-commit]]
type = "write-file"
path = "repo-info.txt"
content = "Repository: {{RepoName}}"
"#;

    ctx.setup_and_install(config);
    ctx.handle_success("pre-commit");

    assert!(ctx.repo.file_exists("repo-info.txt"));

    let content = ctx.repo.read_file("repo-info.txt");
    assert!(content.starts_with("Repository: "));
}

/// Tests that multiple variables from different sources can be extracted and used together.
/// Verifies simultaneous branch and repo_path variable extraction and template rendering.
#[test]
fn template_multiple_variables() {
    let ctx = TestContext::new();

    let config = r#"
extract = [
    "branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)",
    "repo_path:.*/(?P<RepoName>[^/]+)$"
]

[[hooks.pre-commit]]
type = "write-file"
path = "info.txt"
content = "Type: {{Type}}, Ticket: {{Ticket}}, Repo: {{RepoName}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/ABC-456");

    ctx.handle_success("pre-commit");
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
    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "exec"
command = "cmd"
args = ["/C", "echo", "{{Feature}}"]
"#;

    #[cfg(not(windows))]
    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "exec"
command = "echo"
args = ["{{Feature}}"]
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/payment");

    ctx.handle_success("pre-commit");
}

/// Tests that optional variables (branch?) are extracted when pattern matches.
/// Verifies optional extraction syntax works correctly when branch name matches pattern.
#[test]
fn template_optional_variable_present() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Feature: {{Feature}}"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("feature/auth");

    ctx.handle_success("pre-commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "Feature: auth");
}

/// Tests that optional variables (branch?) don't cause failure when pattern doesn't match.
/// Verifies that optional extraction allows hook to proceed even when variable isn't extracted.
#[test]
fn template_optional_variable_missing() {
    let ctx = TestContext::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^.+$"
"#;

    ctx.setup_and_install(config);
    ctx.repo.create_branch("bugfix/issue");

    ctx.handle_success("pre-commit");
}

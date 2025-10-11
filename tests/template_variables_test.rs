mod common;

use common::TestContext;

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

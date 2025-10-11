mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn template_branch_variable_in_message_prefix() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)/(?P<Ticket>[A-Z]+-\\d+)"]

[[hooks.commit-msg]]
type = "message-prefix"
prefix = "{{Type}}: [{{Ticket}}] "
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/PROJ-123");

    repo.write_commit_msg_file("feature: [PROJ-123] add new feature");
    let msg_path = repo.commit_msg_file_path();
    let handle_output = binary.handle(
        "commit-msg",
        repo.path(),
        &[msg_path.to_str().unwrap()],
    );

    assert!(
        handle_output.status.success(),
        "Hook should succeed with templated prefix: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

#[test]
fn template_branch_variable_in_write_file() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "branch-info.txt"
content = "Current feature: {{Feature}}"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/auth-system");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert!(repo.file_exists("branch-info.txt"));
    assert_eq!(
        repo.read_file("branch-info.txt"),
        "Current feature: auth-system"
    );
}

#[test]
fn template_repo_path_variable() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["repo_path:.*/(?P<RepoName>[^/]+)$"]

[[hooks.pre-commit]]
type = "write-file"
path = "repo-info.txt"
content = "Repository: {{RepoName}}"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert!(repo.file_exists("repo-info.txt"));

    let content = repo.read_file("repo-info.txt");
    assert!(content.starts_with("Repository: "));
}

#[test]
fn template_multiple_variables() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/ABC-456");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert!(repo.file_exists("info.txt"));

    let content = repo.read_file("info.txt");
    assert!(content.contains("Type: feature"));
    assert!(content.contains("Ticket: ABC-456"));
    assert!(content.contains("Repo: "));
}

#[test]
fn template_in_exec_command() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/payment");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should succeed with templated exec: {}",
        String::from_utf8_lossy(&handle_output.stderr)
    );
}

#[test]
fn template_optional_variable_present() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "output.txt"
content = "Feature: {{Feature}}"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("feature/auth");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(handle_output.status.success());
    assert_eq!(repo.read_file("output.txt"), "Feature: auth");
}

#[test]
fn template_optional_variable_missing() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "branch-name-regex"
regex = "^.+$"
"#;

    repo.create_config(config);
    repo.create_file("test.txt", "initial");
    let _ = repo.commit("initial");

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/issue");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);

    assert!(
        handle_output.status.success(),
        "Hook should succeed when optional variable is missing"
    );
}

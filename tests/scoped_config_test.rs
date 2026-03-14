mod common;

use common::test_context::TestContext;
use common::{ConfigBuilder, ConfigFormat, FishermanBinary};
// NOTE: Global config tests are not included because the `dirs` crate caches

#[test]
fn repository_and_local_config_merge() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "repo.txt"
content = "repo rule"
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "local.txt"
content = "local rule"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(
        output.status.success(),
        "Both repo and local configs should execute: {}",
        String::from_utf8_lossy(&output.stderr)
    );

    assert!(ctx.repo.file_exists("repo.txt"), "Repo config rule should execute");
    assert!(ctx.repo.file_exists("local.txt"), "Local config rule should execute");
}

#[test]
fn configs_are_concatenated_not_replaced() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "order1.txt"
content = "first"

[[hooks.pre-commit]]
type = "write-file"
path = "order2.txt"
content = "second"
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "order3.txt"
content = "third"

[[hooks.pre-commit]]
type = "write-file"
path = "order4.txt"
content = "fourth"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("order1.txt"));
    assert!(ctx.repo.file_exists("order2.txt"));
    assert!(ctx.repo.file_exists("order3.txt"));
    assert!(ctx.repo.file_exists("order4.txt"));
}

#[test]
fn local_and_repository_both_execute() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "from-repo.txt"
content = "repository config"
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "from-local.txt"
content = "local config"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("from-repo.txt"));
    assert!(ctx.repo.file_exists("from-local.txt"));
}

#[test]
fn different_hooks_in_different_scopes() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "pre-commit-repo.txt"
content = "repo pre-commit"
"#;

    let local_config = r#"
[[hooks.commit-msg]]
type = "message-regex"
regex = ".*"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    assert!(ctx.repo.hook_exists("pre-commit"));
    assert!(ctx.repo.hook_exists("commit-msg"));

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("pre-commit-repo.txt"));
}

#[test]
fn mixed_config_formats_across_scopes() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let yaml_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: yaml-config.txt
      content: from yaml
"#;

    let toml_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "toml-config.txt"
content = "from toml"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository_with_format(ConfigFormat::Yaml, yaml_config)
        .local(toml_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("yaml-config.txt"));
    assert!(ctx.repo.file_exists("toml-config.txt"));
}

#[test]
fn local_config_with_templates() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let local_config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "branch-type.txt"
content = "Type: {{Type}}"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    ctx.repo.create_branch("feature/test");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    let content = ctx.repo.read_file("branch-type.txt");
    assert_eq!(content, "Type: feature");
}

#[test]
fn local_extract_patterns() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let local_config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "feature.txt"
content = "{{Feature}}"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    ctx.repo.create_branch("feature/auth");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    let content = ctx.repo.read_file("feature.txt");
    assert_eq!(content, "auth");
}

#[test]
fn conditional_rules_across_scopes() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "features-only.txt"
content = "feature branch"
when = "Type == \"feature\""
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "always.txt"
content = "always executed"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(ctx.repo.path(), false);

    ctx.repo.create_branch("feature/test");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("features-only.txt"), "Conditional rule should execute");
    assert!(ctx.repo.file_exists("always.txt"), "Unconditional rule should execute");
}

#[test]
fn repository_only_without_global() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "repo-only.txt"
content = "repository only"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("repo-only.txt"));
}

#[test]
fn local_only_without_repository_config() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let local_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "local-only.txt"
content = "local only"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("local-only.txt"));
}

#[test]
fn local_config_yaml_format() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let yaml_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: local-yaml.txt
      content: from local yaml
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .local_with_format(ConfigFormat::Yaml, yaml_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("local-yaml.txt"));

    let content = ctx.repo.read_file("local-yaml.txt");
    assert_eq!(content, "from local yaml");
}

#[test]
fn local_config_json_format() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let json_config = r#"{
  "hooks": {
    "pre-commit": [
      {
        "type": "write-file",
        "path": "local-json.txt",
        "content": "from local json"
      }
    ]
  }
}"#;

    ConfigBuilder::new(&mut ctx.repo)
        .local_with_format(ConfigFormat::Json, json_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());
    assert!(ctx.repo.file_exists("local-json.txt"));

    let content = ctx.repo.read_file("local-json.txt");
    assert_eq!(content, "from local json");
}

#[test]
fn repository_toml_local_yaml() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let toml_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "from-toml.txt"
content = "repository toml"
"#;

    let yaml_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: from-yaml.txt
      content: local yaml
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(toml_config)
        .local_with_format(ConfigFormat::Yaml, yaml_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("from-toml.txt"));
    assert!(ctx.repo.file_exists("from-yaml.txt"));
}

#[test]
fn repository_json_local_yaml() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let json_config = r#"{
  "hooks": {
    "pre-commit": [
      {
        "type": "write-file",
        "path": "from-json.txt",
        "content": "repository json"
      }
    ]
  }
}"#;

    let yaml_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: from-yaml.txt
      content: local yaml
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository_with_format(ConfigFormat::Json, json_config)
        .local_with_format(ConfigFormat::Yaml, yaml_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(output.status.success());

    assert!(ctx.repo.file_exists("from-json.txt"));
    assert!(ctx.repo.file_exists("from-yaml.txt"));
}

#[test]
fn validation_failure_in_any_scope_fails_hook() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = r#"
[[hooks.pre-commit]]
type = "branch-name-prefix"
prefix = "feature/"
"#;

    let local_config = r#"
[[hooks.pre-commit]]
type = "branch-name-suffix"
suffix = "-dev"
"#;

    ConfigBuilder::new(&mut ctx.repo)
        .repository(repo_config)
        .local(local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    ctx.repo.create_branch("bugfix/test-dev");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success(), "Hook should fail when repository validation fails");

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(stderr.contains("feature/"), "Error should mention the failed prefix requirement");
}

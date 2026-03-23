mod common;

use common::test_context::TestContext;
use common::{ConfigBuilder, ConfigFormat, FishermanBinary};
use fisherman_core::{BranchNamePrefixRule, BranchNameSuffixRule, CommitMessageRegexRule, Configuration, Expression, GitHook, WriteFileRule, t};
// NOTE: Global config tests are not included because the `dirs` crate caches

#[test]
fn repository_and_local_config_merge() {
    let binary = FishermanBinary::build();
    let mut ctx = TestContext::new();

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "repo.txt".into(),
                content: "repo rule".into(),
                append: None,
            })
        ]
    );

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "local.txt".into(),
                content: "local rule".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "order1.txt".into(),
                content: "first".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "order2.txt".into(),
                content: "second".into(),
                append: None,
            })
        ]
    );

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "order3.txt".into(),
                content: "third".into(),
                append: None,
            }),
            rule!(WriteFileRule {
                path: "order4.txt".into(),
                content: "fourth".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-repo.txt".into(),
                content: "repository config".into(),
                append: None,
            })
        ]
    );

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-local.txt".into(),
                content: "local config".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "pre-commit-repo.txt".into(),
                content: "repo pre-commit".into(),
                append: None,
            })
        ]
    );

    let local_config = config!(
        GitHook::CommitMsg => [
            rule!(CommitMessageRegexRule {
                when: None,
                expression: ".*".into(),
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
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

    let yaml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "yaml-config.txt".into(),
                content: "from yaml".into(),
                append: None,
            })
        ]
    );

    let toml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "toml-config.txt".into(),
                content: "from toml".into(),
                append: None,
            })
        ]
    );

    let yaml_string = crate::common::configuration::serialize_configuration(&yaml_config, ConfigFormat::Yaml);
    let toml_string = crate::common::configuration::serialize_configuration(&toml_config, ConfigFormat::Toml);

    ConfigBuilder::new(&mut ctx.repo)
        .repository_with_format(ConfigFormat::Yaml, &yaml_string)
        .local_with_format(ConfigFormat::Toml, &toml_string)
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

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "branch-type.txt".into(),
                content: t!("Type: {{Type}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .local_config(&local_config)
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

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "feature.txt".into(),
                content: t!("{{Feature}}"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .local_config(&local_config)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "features-only.txt".into(),
                content: "feature branch".into(),
                append: None,
            }, when = Expression::new("Type == \"feature\""))
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "always.txt".into(),
                content: "always executed".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "repo-only.txt".into(),
                content: "repository only".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
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

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "local-only.txt".into(),
                content: "local only".into(),
                append: None,
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .local_config(&local_config)
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

    let yaml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "local-yaml.txt".into(),
                content: "from local yaml".into(),
                append: None,
            })
        ]
    );

    let yaml_string = crate::common::configuration::serialize_configuration(&yaml_config, ConfigFormat::Yaml);

    ConfigBuilder::new(&mut ctx.repo)
        .local_with_format(ConfigFormat::Yaml, &yaml_string)
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

    let json_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "local-json.txt".into(),
                content: "from local json".into(),
                append: None,
            })
        ]
    );

    let json_string = crate::common::configuration::serialize_configuration(&json_config, ConfigFormat::Json);

    ConfigBuilder::new(&mut ctx.repo)
        .local_with_format(ConfigFormat::Json, &json_string)
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

    let toml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-toml.txt".into(),
                content: "repository toml".into(),
                append: None,
            })
        ]
    );

    let yaml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-yaml.txt".into(),
                content: "local yaml".into(),
                append: None,
            })
        ]
    );

    let toml_string = crate::common::configuration::serialize_configuration(&toml_config, ConfigFormat::Toml);
    let yaml_string = crate::common::configuration::serialize_configuration(&yaml_config, ConfigFormat::Yaml);

    ConfigBuilder::new(&mut ctx.repo)
        .repository_with_format(ConfigFormat::Toml, &toml_string)
        .local_with_format(ConfigFormat::Yaml, &yaml_string)
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

    let json_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-json.txt".into(),
                content: "repository json".into(),
                append: None,
            })
        ]
    );

    let yaml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: "from-yaml.txt".into(),
                content: "local yaml".into(),
                append: None,
            })
        ]
    );

    let json_string = crate::common::configuration::serialize_configuration(&json_config, ConfigFormat::Json);
    let yaml_string = crate::common::configuration::serialize_configuration(&yaml_config, ConfigFormat::Yaml);

    ConfigBuilder::new(&mut ctx.repo)
        .repository_with_format(ConfigFormat::Json, &json_string)
        .local_with_format(ConfigFormat::Yaml, &yaml_string)
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

    let repo_config = config!(
        GitHook::PreCommit => [
            rule!(BranchNamePrefixRule {
                prefix: "feature/".into(),
            })
        ]
    );

    let local_config = config!(
        GitHook::PreCommit => [
            rule!(BranchNameSuffixRule {
                suffix: "-dev".into(),
            })
        ]
    );

    ConfigBuilder::new(&mut ctx.repo)
        .repository_config(&repo_config)
        .local_config(&local_config)
        .build();

    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);
    binary.install(ctx.repo.path(), false);

    ctx.repo.create_branch("bugfix/test-dev");

    let output = ctx.git_commit_allow_empty("test commit");
    assert!(!output.status.success(), "Hook should fail when repository validation fails");

    let stderr = String::from_utf8_lossy(&output.stderr);
    assert!(stderr.contains("feature/"), "Error should mention the failed prefix requirement");
}

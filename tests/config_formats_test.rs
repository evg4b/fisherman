mod common;

use crate::common::ConfigFormat;
use crate::common::configuration::serialize_configuration;
use common::test_context::TestContext;
use fisherman_core::{t, Configuration, GitHook, WriteFileRule};

#[test]
fn yaml_config_format() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: t!("yaml-test.txt"),
                content: t!("YAML config works"),
                append: None,
            })
        ]
    );

    let yaml_string = serialize_configuration(&config, ConfigFormat::Yaml);
    ctx.repo.create_yaml_config(&yaml_string);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = ctx.binary.install(ctx.repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with YAML config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("yaml-test.txt"));
    assert_eq!(ctx.repo.read_file("yaml-test.txt"), "YAML config works");
}

#[test]
fn json_config_format() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: t!("json-test.txt"),
                content: t!("JSON config works"),
                append: None,
            })
        ]
    );

    let json_string = serialize_configuration(&config, ConfigFormat::Json);
    ctx.repo.create_json_config(&json_string);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = ctx.binary.install(ctx.repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with JSON config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("json-test.txt"));
    assert_eq!(ctx.repo.read_file("json-test.txt"), "JSON config works");
}

#[test]
fn yaml_with_templates() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: t!("output.txt"),
                content: t!("Branch type: {{Type}}"),
                append: None,
            })
        ],
        extract = vec![String::from("branch:^(?P<Type>feature|bugfix)")]
    );

    let yaml_string = serialize_configuration(&config, ConfigFormat::Yaml);
    ctx.repo.create_yaml_config(&yaml_string);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    ctx.binary.install(ctx.repo.path(), false);
    ctx.repo.create_branch("feature/test");

    ctx.git_commit_allow_empty_success("test commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "Branch type: feature");
}

#[test]
fn json_with_conditional() {
    let ctx = TestContext::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: t!("conditional.txt"),
                    content: t!("Feature branch"),
                    append: None,
                },
                when = fisherman_core::Expression::new("Type == \"feature\"")
            )
        ],
        extract = vec![String::from("branch:^(?P<Type>feature|bugfix)")]
    );

    let json_string = serialize_configuration(&config, ConfigFormat::Json);
    ctx.repo.create_json_config(&json_string);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    ctx.binary.install(ctx.repo.path(), false);
    ctx.repo.create_branch("feature/test");

    ctx.git_commit_allow_empty_success("test commit");
    assert!(ctx.repo.file_exists("conditional.txt"));
}

#[test]
fn multiple_config_formats_error() {
    let ctx = TestContext::new();

    let toml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: t!("toml.txt"),
                content: t!("TOML"),
                append: None,
            })
        ]
    );

    let yaml_config = config!(
        GitHook::PreCommit => [
            rule!(WriteFileRule {
                path: t!("yaml.txt"),
                content: t!("YAML"),
                append: None,
            })
        ]
    );

    let toml_string = serialize_configuration(&toml_config, ConfigFormat::Toml);
    let yaml_string = serialize_configuration(&yaml_config, ConfigFormat::Yaml);

    ctx.repo.create_config(&toml_string, ConfigFormat::Toml);
    ctx.repo.create_yaml_config(&yaml_string);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = ctx.binary.install(ctx.repo.path(), false);

    if install_output.status.success() {
        let handle_output = ctx.git_commit_allow_empty("test commit");
        assert!(
            !handle_output.status.success(),
            "Should fail when multiple config formats are present"
        );
    } else {
        let stderr = String::from_utf8_lossy(&install_output.stderr);
        assert!(
            stderr.contains("Multiple config files"),
            "Error should mention multiple config files"
        );
    }
}

mod common;

use common::test_context::TestContext;

#[test]
fn yaml_config_format() {
    let ctx = TestContext::new();

    let config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: yaml-test.txt
      content: YAML config works
"#;

    ctx.repo.create_yaml_config(config);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = ctx.binary.install(ctx.repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with YAML config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("yaml-test.txt"));
    assert_eq!(ctx.repo.read_file("yaml-test.txt"), "YAML config works");
}

#[test]
fn json_config_format() {
    let ctx = TestContext::new();

    let config = r#"
{
  "hooks": {
    "pre-commit": [
      {
        "type": "write-file",
        "path": "json-test.txt",
        "content": "JSON config works"
      }
    ]
  }
}
"#;

    ctx.repo.create_json_config(config);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = ctx.binary.install(ctx.repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with JSON config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("json-test.txt"));
    assert_eq!(ctx.repo.read_file("json-test.txt"), "JSON config works");
}

#[test]
fn yaml_with_templates() {
    let ctx = TestContext::new();

    let config = r#"
extract:
  - "branch:^(?P<Type>feature|bugfix)"
hooks:
  pre-commit:
    - type: write-file
      path: output.txt
      content: "Branch type: {{Type}}"
"#;

    ctx.repo.create_yaml_config(config);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    ctx.binary.install(ctx.repo.path(), false);
    ctx.repo.create_branch("feature/test");

    ctx.handle_success("pre-commit");
    assert_eq!(ctx.repo.read_file("output.txt"), "Branch type: feature");
}

#[test]
fn json_with_conditional() {
    let ctx = TestContext::new();

    let config = r#"
{
  "extract": ["branch:^(?P<Type>feature|bugfix)"],
  "hooks": {
    "pre-commit": [
      {
        "type": "write-file",
        "path": "conditional.txt",
        "content": "Feature branch",
        "when": "Type == \"feature\""
      }
    ]
  }
}
"#;

    ctx.repo.create_json_config(config);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    ctx.binary.install(ctx.repo.path(), false);
    ctx.repo.create_branch("feature/test");

    ctx.handle_success("pre-commit");
    assert!(ctx.repo.file_exists("conditional.txt"));
}

#[test]
fn multiple_config_formats_error() {
    let ctx = TestContext::new();

    // Create both TOML and YAML configs
    let toml_config = r#"
[[hooks.pre-commit]]
type = "write-file"
path = "toml.txt"
content = "TOML"
"#;

    let yaml_config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: yaml.txt
      content: YAML
"#;

    ctx.repo.create_config(toml_config);
    ctx.repo.create_yaml_config(yaml_config);
    ctx.repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    // Installation or handling should fail with multiple config formats
    let install_output = ctx.binary.install(ctx.repo.path(), false);

    if install_output.status.success() {
        let handle_output = ctx.handle("pre-commit");
        assert!(
            !handle_output.status.success(),
            "Should fail when multiple config formats are present"
        );
    } else {
        // Installation itself failed, which is also acceptable
        let stderr = String::from_utf8_lossy(&install_output.stderr);
        assert!(
            stderr.contains("Multiple config files"),
            "Error should mention multiple config files"
        );
    }
}

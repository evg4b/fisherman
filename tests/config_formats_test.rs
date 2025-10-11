mod common;

use common::{FishermanBinary, GitTestRepo};

#[test]
fn yaml_config_format() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
hooks:
  pre-commit:
    - type: write-file
      path: yaml-test.txt
      content: YAML config works
"#;

    repo.create_yaml_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with YAML config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(handle_output.status.success());
    assert!(repo.file_exists("yaml-test.txt"));
    assert_eq!(repo.read_file("yaml-test.txt"), "YAML config works");
}

#[test]
fn json_config_format() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_json_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    let install_output = binary.install(repo.path(), false);
    assert!(
        install_output.status.success(),
        "Installation should succeed with JSON config: {}",
        String::from_utf8_lossy(&install_output.stderr)
    );

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(handle_output.status.success());
    assert!(repo.file_exists("json-test.txt"));
    assert_eq!(repo.read_file("json-test.txt"), "JSON config works");
}

#[test]
fn yaml_with_templates() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract:
  - "branch:^(?P<Type>feature|bugfix)"
hooks:
  pre-commit:
    - type: write-file
      path: output.txt
      content: "Branch type: {{Type}}"
"#;

    repo.create_yaml_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);
    repo.create_branch("feature/test");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(handle_output.status.success());
    assert_eq!(repo.read_file("output.txt"), "Branch type: feature");
}

#[test]
fn json_with_conditional() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_json_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);
    repo.create_branch("feature/test");

    let handle_output = binary.handle("pre-commit", repo.path(), &[]);
    assert!(handle_output.status.success());
    assert!(repo.file_exists("conditional.txt"));
}

#[test]
fn multiple_config_formats_error() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

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

    repo.create_config(toml_config);
    repo.create_yaml_config(yaml_config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    // Installation or handling should fail with multiple config formats
    let install_output = binary.install(repo.path(), false);

    if install_output.status.success() {
        let handle_output = binary.handle("pre-commit", repo.path(), &[]);
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

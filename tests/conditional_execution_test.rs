mod common;

use crate::common::ConfigFormat;
use common::{configuration::serialize_configuration, FishermanBinary, GitTestRepo};
use core::configuration::Configuration;
use core::hooks::GitHook;
use core::rules::RuleParams;

#[test]
fn when_condition_true_executes_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("executed.txt"),
                    content: String::from("Rule executed"),
                    append: None,
                },
                when = String::from("Type == \"feature\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(
        repo.file_exists("executed.txt"),
        "File should be created when condition is true"
    );
}

#[test]
fn when_condition_false_skips_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("executed.txt"),
                    content: String::from("Rule executed"),
                    append: None,
                },
                when = String::from("Type == \"feature\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(
        !repo.file_exists("executed.txt"),
        "File should not be created when condition is false"
    );
}

#[test]
fn when_condition_with_is_def_var_defined() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("feature.txt"),
                    content: String::from("Feature: {{Feature}}"),
                    append: None,
                },
                when = String::from("is_def_var(\"Feature\")")
            )
        ],
        extract = vec![
            String::from("branch:^feature/(?P<Feature>[a-z-]+)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/auth");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("feature.txt"));
}

#[test]
fn when_condition_with_is_def_var_undefined() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("feature.txt"),
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

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(
        !repo.file_exists("feature.txt"),
        "File should not be created when variable is undefined"
    );
}

#[test]
fn when_condition_complex_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("urgent.txt"),
                    content: String::from("Urgent feature"),
                    append: None,
                },
                when = String::from("Type == \"feature\" && Priority == \"high\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Priority>high|low)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/high");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(
        repo.file_exists("urgent.txt"),
        "File should be created when complex condition is true"
    );
}

#[test]
fn when_condition_complex_expression_false() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("urgent.txt"),
                    content: String::from("Urgent feature"),
                    append: None,
                },
                when = String::from("Type == \"feature\" && Priority == \"high\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)/(?P<Priority>high|low)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/low");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(!repo.file_exists("urgent.txt"));
}

#[test]
fn when_condition_or_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("production.txt"),
                    content: String::from("Production change"),
                    append: None,
                },
                when = String::from("Type == \"hotfix\" || Type == \"bugfix\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix|hotfix)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("hotfix/urgent");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("production.txt"));
}

#[test]
fn when_condition_not_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("not-feature.txt"),
                    content: String::from("Not a feature"),
                    append: None,
                },
                when = String::from("Type != \"feature\"")
            )
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("not-feature.txt"));
}

#[test]
fn when_condition_multiple_rules_selective_execution() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = config!(
        GitHook::PreCommit => [
            rule!(
                WriteFileRule {
                    path: String::from("feature.txt"),
                    content: String::from("Feature branch"),
                    append: None,
                },
                when = String::from("Type == \"feature\"")
            ),
            rule!(
                WriteFileRule {
                    path: String::from("bugfix.txt"),
                    content: String::from("Bugfix branch"),
                    append: None,
                },
                when = String::from("Type == \"bugfix\"")
            ),
            rule!(WriteFileRule {
                path: String::from("always.txt"),
                content: String::from("Always executed"),
                append: None,
            })
        ],
        extract = vec![
            String::from("branch:^(?P<Type>feature|bugfix)"),
        ]
    );

    let config_string = serialize_configuration(&config, ConfigFormat::Toml);
    repo.create_config(&config_string, ConfigFormat::Toml);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("feature.txt"));
    assert!(!repo.file_exists("bugfix.txt"));
    assert!(repo.file_exists("always.txt"));
}

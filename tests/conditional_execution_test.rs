mod common;

use common::{FishermanBinary, GitTestRepo};

/// Tests that a rule with a true 'when' condition executes successfully.
/// Verifies conditional execution based on extracted variables.
#[test]
fn when_condition_true_executes_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "executed.txt"
content = "Rule executed"
when = "Type == \"feature\""
"#;

    repo.create_config(config);
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

/// Tests that a rule with a false 'when' condition is skipped and doesn't execute.
/// Verifies rules are properly filtered based on conditions.
#[test]
fn when_condition_false_skips_rule() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "executed.txt"
content = "Rule executed"
when = "Type == \"feature\""
"#;

    repo.create_config(config);
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

/// Tests that is_def_var() function works correctly when a variable is defined.
/// Verifies conditional execution can check for variable existence.
#[test]
fn when_condition_with_is_def_var_defined() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "feature.txt"
content = "Feature: {{Feature}}"
when = "is_def_var(\"Feature\")"
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/auth");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("feature.txt"));
}

/// Tests that is_def_var() returns false when a variable is undefined (optional extraction).
/// Verifies rules can gracefully handle missing variables.
#[test]
fn when_condition_with_is_def_var_undefined() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch?:^feature/(?P<Feature>[a-z-]+)"]

[[hooks.pre-commit]]
type = "write-file"
path = "feature.txt"
content = "Feature: {{Feature}}"
when = "is_def_var(\"Feature\")"
"#;

    repo.create_config(config);
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

/// Tests that complex AND expressions in 'when' conditions work correctly.
/// Verifies multiple conditions can be combined with && operator.
#[test]
fn when_condition_complex_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)/(?P<Priority>high|low)"]

[[hooks.pre-commit]]
type = "write-file"
path = "urgent.txt"
content = "Urgent feature"
when = "Type == \"feature\" && Priority == \"high\""
"#;

    repo.create_config(config);
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

/// Tests that complex AND expressions correctly evaluate to false when one condition fails.
/// Verifies && operator requires all conditions to be true.
#[test]
fn when_condition_complex_expression_false() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)/(?P<Priority>high|low)"]

[[hooks.pre-commit]]
type = "write-file"
path = "urgent.txt"
content = "Urgent feature"
when = "Type == \"feature\" && Priority == \"high\""
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/low");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(!repo.file_exists("urgent.txt"));
}

/// Tests that OR expressions in 'when' conditions work correctly.
/// Verifies conditions can be combined with || operator.
#[test]
fn when_condition_or_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix|hotfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "production.txt"
content = "Production change"
when = "Type == \"hotfix\" || Type == \"bugfix\""
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("hotfix/urgent");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("production.txt"));
}

/// Tests that NOT (!=) expressions in 'when' conditions work correctly.
/// Verifies inequality comparisons can be used in conditions.
#[test]
fn when_condition_not_expression() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "not-feature.txt"
content = "Not a feature"
when = "Type != \"feature\""
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("bugfix/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("not-feature.txt"));
}

/// Tests that multiple rules with different conditions execute selectively.
/// Verifies each rule's condition is evaluated independently.
#[test]
fn when_condition_multiple_rules_selective_execution() {
    let binary = FishermanBinary::build();
    let repo = GitTestRepo::new();

    let config = r#"
extract = ["branch:^(?P<Type>feature|bugfix)"]

[[hooks.pre-commit]]
type = "write-file"
path = "feature.txt"
content = "Feature branch"
when = "Type == \"feature\""

[[hooks.pre-commit]]
type = "write-file"
path = "bugfix.txt"
content = "Bugfix branch"
when = "Type == \"bugfix\""

[[hooks.pre-commit]]
type = "write-file"
path = "always.txt"
content = "Always executed"
"#;

    repo.create_config(config);
    repo.git_history(&[("initial", &[("test.txt", "initial")])]);

    binary.install(repo.path(), false);

    repo.create_branch("feature/test");

    let handle_output = repo.commit_with_hooks_allow_empty("test commit");

    assert!(handle_output.status.success());
    assert!(repo.file_exists("feature.txt"));
    assert!(!repo.file_exists("bugfix.txt"));
    assert!(repo.file_exists("always.txt"));
}

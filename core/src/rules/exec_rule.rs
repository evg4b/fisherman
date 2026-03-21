use crate::context::Context;
use crate::extract_vars;
use crate::rules::rule::{ConditionalRule, Rule, RuleResult};
use crate::scripting::Expression;
use crate::templates::{replace_in_hashmap, replace_in_vec};
use anyhow::Result;
use rules_derive::ConditionalRule as ConditionalRuleDerive;
use std::collections::HashMap;
use std::env;
use std::process::Command;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct ExecRule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    pub command: String,
    pub args: Option<Args>,
    pub env: Option<Env>,
}

impl std::fmt::Display for ExecRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let args = self.args.clone().unwrap_or_default();
        let mut command = self.command.clone();
        if !args.is_empty() {
            command.push(' ');
            command.push_str(&args.join(" "));
        }
        write!(f, "Execute command: {}", command)
    }
}

#[typetag::serde(name = "exec")]
impl Rule for ExecRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.when.is_some() && !self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: "exec".into(),
            });
        }

        let variables = extract_vars!(self, ctx)?;
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(&self.env.clone().unwrap_or_default(), &variables)?);

        let output = Command::new(self.command.clone())
            .args(replace_in_vec(&self.args.clone().unwrap_or(vec![]), &variables)?)
            .envs(env_map)
            .output()?;

        match output.status.success() {
            true => Ok(RuleResult::Success {
                name: "exec".into(),
                output: Some(String::from_utf8_lossy(&output.stdout).to_string()),
            }),
            false => Ok(RuleResult::Failure {
                name: "exec".into(),
                message: String::from_utf8_lossy(&output.stderr).to_string(),
            }),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::context::MockContext;
    use assert2::assert;
    use std::collections::HashMap;

    fn mock_ctx_with_vars(vars: HashMap<String, String>) -> MockContext {
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));
        ctx
    }

    fn mock_ctx() -> MockContext {
        mock_ctx_with_vars(HashMap::new())
    }

    #[test]
    fn test_exec_rule() {
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["hello".into()]),
            env: None,
            when: None,
            extract: None,
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        assert_eq!(output.unwrap(), "hello\n");
    }

    #[test]
    fn test_exec_rule_new() {
        let rule = ExecRule {
            when: None,
            extract: None,
            command: "echo".into(),
            args: Some(vec!["hello".into()]),
            env: None,
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        assert_eq!(output.unwrap(), "hello\n");
    }

    #[test]
    fn test_exec_rule_with_env() {
        let mut env = HashMap::new();
        env.insert("HELLO".into(), "world".into());

        let rule = ExecRule {
            when: None,
            extract: None,
            command: "printenv".into(),
            args: None,
            env: Some(env),
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        assert!(output.unwrap().contains("HELLO=world"));
    }

    #[test]
    fn test_exec_rule_with_variables() {
        let mut vars = HashMap::new();
        vars.insert("HELLO".into(), "world".into());

        let rule = ExecRule {
            when: None,
            extract: None,
            command: "echo".into(),
            args: Some(vec!["hello".into(), "{{HELLO}}".into()]),
            env: None,
        };

        let result = rule.check(&mock_ctx_with_vars(vars)).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        assert_eq!(output.unwrap(), "hello world\n");
    }

    #[test]
    fn test_exec_rule_failure_on_missing_file() {
        let rule = ExecRule {
            when: None,
            extract: None,
            command: "cat".into(),
            args: Some(vec!["./unknown.txt".into()]),
            env: None,
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "exec");
        assert!(message.contains("No such file or directory"));
    }

    #[test]
    fn test_return_error() {
        let rule = ExecRule {
            when: None,
            extract: None,
            command: "XXXXXXXXXXXX".into(),
            args: None,
            env: None,
        };

        let result = rule.check(&mock_ctx()).err().unwrap();

        assert_eq!(result.to_string(), "No such file or directory (os error 2)");
    }

    #[test]
    fn test_exec_rule_env_template_error() {
        let mut env = HashMap::new();
        env.insert("VAR".into(), "{{missing}}".into());

        let rule = ExecRule {
            when: None,
            extract: None,
            command: "echo".into(),
            args: Some(vec!["{{VAR}}".into()]),
            env: Some(env),
        };

        let result = rule.check(&mock_ctx());
        assert!(result.is_err());
    }

    #[test]
    fn test_exec_rule_args_template_error() {
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["{{VAR}}".into()]),
            env: None,
            when: None,
            extract: None,
        };

        let result = rule.check(&mock_ctx());
        assert!(result.is_err());
    }

    #[test]
    fn test_display_with_args() {
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["hello".into(), "world".into()]),
            env: None,
            when: None,
            extract: None,
        };
        assert_eq!(format!("{}", rule), "Execute command: echo hello world");
    }

    #[test]
    fn test_display_without_args() {
        let rule = ExecRule {
            command: "ls".into(),
            args: None,
            env: None,
            when: None,
            extract: None,
        };
        assert_eq!(format!("{}", rule), "Execute command: ls");
    }
}

use crate::context::Context;
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
    pub command: String,
    pub args: Option<Args>,
    pub env: Option<Env>,
}

#[typetag::serde(name = "exec")]
impl Rule for ExecRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: "exec".into(),
            });
        }

        let variables = HashMap::new();
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(&self.env.clone().unwrap_or(HashMap::new()), &variables)?);

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

    #[test]
    fn test_exec_rule() {
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["hello".into()]),
            env: None,
            when: None,
        };

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert_eq!(output.unwrap(), "hello\n");
    }

    #[test]
    fn test_exec_rule_new() {
        let rule = ExecRule {
            when: None,
            command: "echo".into(),
            args: Some(vec!["hello".into()]),
            env: None,
        };

        let result = rule.check(&MockContext::new()).unwrap();
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
            command: "printenv".into(),
            args: None,
            env: Some(env),
        };

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert!(output.unwrap().contains("HELLO=world"));
    }

    #[test]
    fn test_exec_rule_with_variables() {
        // let mut variables = HashMap::new();
        // variables.insert("HELLO".into(), "world".into());

        let rule = ExecRule {
            when: None,
            command: "echo".into(),
            args: Some(vec!["hello".into(), "{{HELLO}}".into()]),
            env: None,
        };

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert_eq!(output.unwrap(), "hello world\n");
    }

    #[test]
    fn test_exec_rule_failure_on_missing_file() {
        // let mut variables = HashMap::new();
        // variables.insert("HELLO".into(), "world".into());

        let rule = ExecRule {
            when: None,
            command: "cat".into(),
            args: Some(vec!["./unknown.txt".into()]),
            env: None,
        };

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "test");
        assert_eq!(message, "cat: ./unknown.txt: No such file or directory\n");
    }

    #[test]
    fn test_return_error() {
        let rule = ExecRule {
            when: None,
            command: "XXXXXXXXXXXX".into(),
            args: None,
            env: None,
        };

        let result = rule.check(&MockContext::new()).err().unwrap();

        assert_eq!(result.to_string(), "No such file or directory (os error 2)");
    }

    #[test]
    fn test_exec_rule_env_template_error() {
        let mut env = HashMap::new();
        env.insert("VAR".into(), "{{missing}}".into());

        let rule = ExecRule {
            when: None,
            command: "echo".into(),
            args: Some(vec!["{{VAR}}".into()]),
            env: Some(env),
        };

        let result = rule.check(&MockContext::new());
        assert!(result.is_err());
    }

    #[test]
    fn test_exec_rule_args_template_error() {
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["{{VAR}}".into()]),
            env: None,
            when: None,
        };

        let result = rule.check(&MockContext::new());
        assert!(result.is_err());
    }
}

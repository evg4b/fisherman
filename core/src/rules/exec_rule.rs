use crate::context::Context;
use crate::rules::compiled_rule::{CompiledRule, RuleResultOld};
use crate::rules::rule::{Rule, RuleResult, ConditionalRule};
use crate::templates::{replace_in_hashmap, replace_in_vec};
use crate::scripting::Expression;
use anyhow::Result;
use std::collections::HashMap;
use std::env;
use std::process::Command;
use rules_derive::ConditionalRule as ConditionalRuleDerive;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

#[derive(Debug)]
pub struct ExecRuleOld {
    name: String,
    command: String,
    args: Args,
    env: Env,
    variables: HashMap<String, String>,
}

impl ExecRuleOld {
    pub fn new(
        name: String,
        command: String,
        args: Args,
        env: Env,
        variables: HashMap<String, String>,
    ) -> Self {
        Self {
            name,
            command,
            args,
            env,
            variables,
        }
    }
}

#[derive(Debug, serde::Serialize, serde::Deserialize, ConditionalRuleDerive)]
pub struct ExecRule {
    pub when: Option<Expression>,
    pub command: String,
    pub args: Args,
    pub env: Env,
}

impl ExecRule {
    pub fn new(command: String, args: Args, env: Env) -> Self {
        Self { when: None, command, args, env }
    }
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
        env_map.extend(replace_in_hashmap(&self.env, &variables)?);

        let output = Command::new(self.command.clone())
            .args(replace_in_vec(&self.args, &variables)?)
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

impl CompiledRule for ExecRuleOld {
    fn is_sequential(&self) -> bool {
        false
    }

    fn check(&self, _: &dyn Context) -> Result<RuleResultOld> {
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(&self.env, &self.variables)?);

        let output = Command::new(self.command.clone())
            .args(replace_in_vec(&self.args, &self.variables)?)
            .envs(env_map)
            .output()?;

        match output.status.success() {
            true => Ok(RuleResultOld::Success {
                name: self.name.clone(),
                output: Some(String::from_utf8_lossy(&output.stdout).to_string()),
            }),
            false => Ok(RuleResultOld::Failure {
                name: self.name.clone(),
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
        let rule = ExecRuleOld::new(
            "test".into(),
            "echo".into(),
            vec!["hello".into()],
            HashMap::new(),
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert_eq!(output.unwrap(), "hello\n");
    }

    #[test]
    fn test_exec_rule_new() {
        let rule = ExecRule { when: None, command: "echo".into(), args: vec!["hello".into()], env: HashMap::new() };

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

        let rule = ExecRuleOld::new(
            "test".into(),
            "printenv".into(),
            vec![],
            env,
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert!(output.unwrap().contains("HELLO=world"));
    }

    #[test]
    fn test_exec_rule_with_variables() {
        let mut variables = HashMap::new();
        variables.insert("HELLO".into(), "world".into());

        let rule = ExecRuleOld::new(
            "test".into(),
            "echo".into(),
            vec!["hello".into(), "{{HELLO}}".into()],
            HashMap::new(),
            variables,
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "test");
        assert_eq!(output.unwrap(), "hello world\n");
    }

    #[test]
    fn test_exec_rule_failure_on_missing_file() {
        let mut variables = HashMap::new();
        variables.insert("HELLO".into(), "world".into());

        let rule = ExecRuleOld::new(
            "test".into(),
            "cat".into(),
            vec!["./unknown.txt".into()],
            HashMap::new(),
            variables,
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResultOld::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "test");
        assert_eq!(message, "cat: ./unknown.txt: No such file or directory\n");
    }

    #[test]
    fn test_return_error() {
        let rule = ExecRuleOld::new(
            "test".into(),
            "XXXXXXXXXXXX".into(),
            vec![],
            HashMap::new(),
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).err().unwrap();

        assert_eq!(result.to_string(), "No such file or directory (os error 2)");
    }

    #[test]
    fn test_is_sequential() {
        let rule = ExecRuleOld::new(
            "test".into(),
            "echo".into(),
            vec!["hello".into()],
            HashMap::new(),
            HashMap::new(),
        );
        assert!(!rule.is_sequential());
    }

    #[test]
    fn test_exec_rule_env_template_error() {
        let mut env = HashMap::new();
        env.insert("VAR".into(), "{{missing}}".into());

        let rule = ExecRuleOld::new("test".into(), "echo".into(), vec![], env, HashMap::new());

        let result = rule.check(&MockContext::new());
        assert!(result.is_err());
    }

    #[test]
    fn test_exec_rule_args_template_error() {
        let rule = ExecRuleOld::new(
            "test".into(),
            "echo".into(),
            vec!["{{missing}}".into()],
            HashMap::new(),
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new());
        assert!(result.is_err());
    }
}

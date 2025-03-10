use crate::rules::compiled_rule::{CompiledRule, RuleResult};
use crate::templates::{replace_in_hashmap, replace_in_vac};
use anyhow::Result;
use std::collections::HashMap;
use std::env;
use std::process::Command;
use crate::context::Context;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

#[derive(Debug)]
pub struct ExecRule {
    name: String,
    command: String,
    args: Args,
    env: Env,
    variables: HashMap<String, String>,
}

impl ExecRule {
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

impl CompiledRule for ExecRule {
    fn check(&self, _: &dyn Context) -> Result<RuleResult> {
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(&self.env, &self.variables)?);

        let output = Command::new(self.command.clone())
            .args(replace_in_vac(&self.args, &self.variables)?)
            .envs(env_map)
            .output()?;

        match output.status.success() {
            true => Ok(RuleResult::Success {
                name: self.name.clone(),
                output: String::from_utf8_lossy(&output.stdout).to_string(),
            }),
            false => Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: String::from_utf8_lossy(&output.stderr).to_string(),
            }),
        }
    }
}

#[cfg(test)]
mod test {
    use super::*;
    use assertor::{assert_that, EqualityAssertion, StringAssertion};
    use std::collections::HashMap;
    use crate::context::MockContext;

    #[test]
    fn test_exec_rule() {
        let rule = ExecRule::new(
            "test".into(),
            "echo".into(),
            vec!["hello".into()],
            HashMap::new(),
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result)
        };

        assert_that!(name).is_equal_to("test".to_string());
        assert_that!(output).is_equal_to("hello\n".to_string());
    }

    #[test]
    fn test_exec_rule_with_env() {
        let mut env = HashMap::new();
        env.insert("HELLO".into(), "world".into());

        let rule = ExecRule::new(
            "test".into(),
            "printenv".into(),
            vec![],
            env,
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result)
        };

        assert_that!(name).is_equal_to("test".to_string());
        assert_that!(output).contains("HELLO=world".to_string());
    }

    #[test]
    fn test_exec_rule_with_variables() {
        let mut variables = HashMap::new();
        variables.insert("HELLO".into(), "world".into());

        let rule = ExecRule::new(
            "test".into(),
            "echo".into(),
            vec!["hello".into(), "{{HELLO}}".into()],
            HashMap::new(),
            variables,
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Expected Success, but got {:?}", result)
        };

        assert_that!(name).is_equal_to("test".to_string());
        assert_that!(output).is_equal_to("hello world\n".to_string());
    }

    #[test]
    fn test_exec_rule_with_variable222() {
        let mut variables = HashMap::new();
        variables.insert("HELLO".into(), "world".into());

        let rule = ExecRule::new(
            "test".into(),
            "cat".into(),
            vec!["./unknown.txt".into()],
            HashMap::new(),
            variables,
        );

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            panic!("Expected Success, but got {:?}", result)
        };

        assert_that!(name).is_equal_to("test".to_string());
        assert_that!(message)
            .is_equal_to("cat: ./unknown.txt: No such file or directory\n".to_string());
    }

    #[test]
    fn test_return_error() {
        let rule = ExecRule::new(
            "test".into(),
            "XXXXXXXXXXXX".into(),
            vec![],
            HashMap::new(),
            HashMap::new(),
        );

        let result = rule.check(&MockContext::new()).err().unwrap();

        assert_that!(result.to_string())
            .is_equal_to("No such file or directory (os error 2)".to_string());
    }
}

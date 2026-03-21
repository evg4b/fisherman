use crate::context::Context;
use crate::rules::rule::{Rule, RuleResult};
use crate::rules::{CompiledRule, RuleResultOld};
use crate::templates::TemplateString;
use anyhow::Result;
use run_script::{run, ScriptOptions};
use std::collections::HashMap;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct ShellScriptRule {
    pub script: String,
    pub env: Option<HashMap<String, String>>,
}

pub struct ShellScript {
    name: String,
    script: TemplateString,
    env: HashMap<String, String>,
    variables: HashMap<String, String>,
}

impl ShellScriptRule {
    pub fn new(script: String, env: Option<HashMap<String, String>>) -> Self {
        Self { script, env }
    }
}

#[typetag::serde(name = "shell")]
impl Rule for ShellScriptRule {
    fn check(&self, _ctx: &dyn Context) -> Result<RuleResult> {
        let mut options = ScriptOptions::new();
        options.env_vars = self.env.clone();

        let args = vec![];
        let (code, output, _) = run(self.script.as_str(), &args, &options)?;

        if code != 0 {
            return Ok(RuleResult::Failure {
                name: "shell".into(),
                message: format!("exit code: {}", code),
            });
        }

        Ok(RuleResult::Success {
            name: "shell".into(),
            output: Some(output),
        })
    }
}

impl ShellScript {
    pub fn new(
        name: String,
        script: TemplateString,
        env: HashMap<String, String>,
        variables: HashMap<String, String>,
    ) -> Self {
        Self {
            name,
            script,
            env,
            variables,
        }
    }
}

impl CompiledRule for ShellScript {
    fn is_sequential(&self) -> bool {
        false
    }

    fn check(&self, _ctx: &dyn Context) -> Result<RuleResultOld> {
        let mut options = ScriptOptions::new();
        options.env_vars = Some(self.env.clone());

        let args = vec![];
        let (code, output, _) = run(
            self.script.compile(&self.variables)?.as_str(),
            &args,
            &options,
        )?;

        if code != 0 {
            return Ok(RuleResultOld::Failure {
                name: self.name.clone(),
                message: format!("exit code: {}", code),
            });
        }

        Ok(RuleResultOld::Success {
            name: self.name.clone(),
            output: Some(output),
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::context::MockContext;
    use crate::rules::rule::{Rule, RuleResult};
    use crate::rules::shell_script::{ShellScript, ShellScriptRule};
    use crate::rules::CompiledRule;
    use crate::rules::RuleResultOld;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_shell_script() {
        let script = ShellScript::new(
            "Test".to_string(),
            t!("echo 'Test'"),
            HashMap::new(),
            HashMap::new(),
        );

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_shell_script_failure() {
        let script = ShellScript::new(
            "Test".to_string(),
            t!("exit 1"),
            HashMap::new(),
            HashMap::new(),
        );

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResultOld::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "Test");
        assert_eq!(message, "exit code: 1");
    }

    #[test]
    fn test_shell_script_with_variables() {
        let mut variables = HashMap::new();
        variables.insert("name".to_string(), "Test".to_string());

        let script = ShellScript::new(
            "Test".to_string(),
            t!("echo 'Hello {{name}}'"),
            HashMap::new(),
            variables,
        );

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Hello Test\n");
    }

    #[test]
    fn test_shell_script_with_env() {
        let mut env = HashMap::new();
        env.insert("TEST".to_string(), "Test".to_string());

        let script = ShellScript::new("Test".to_string(), t!("echo $TEST"), env, HashMap::new());

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResultOld::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_is_sequential() {
        let script = ShellScript::new(
            "Test".to_string(),
            t!("echo 'Test'"),
            HashMap::new(),
            HashMap::new(),
        );
        assert!(!script.is_sequential());
    }

    #[test]
    fn test_shell_script_rule_new() {
        let rule = ShellScriptRule::new("echo 'Test'".to_string(), None);

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_shell_script_template_error() {
        let script = ShellScript::new(
            "Test".to_string(),
            t!("echo '{{missing}}'"),
            HashMap::new(),
            HashMap::new(),
        );

        let result = script.check(&MockContext::new());
        assert!(result.is_err());
    }
}

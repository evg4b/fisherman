use crate::context::Context;
use crate::rules::{ExecutionMode, Rule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use run_script::{run, ScriptOptions};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize)]
pub struct ShellScriptRule {
    pub script: TemplateString,
    pub env: Option<HashMap<String, String>>,
}

impl std::fmt::Display for ShellScriptRule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Run shell script: {}", self.script)
    }
}

static SHELL_SCRIPT_NAME: &str = "shell";

#[typetag::serde(name = "shell")]
impl Rule for ShellScriptRule {
    fn execution_mode(&self) -> ExecutionMode {
        ExecutionMode::Async
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let mut options = ScriptOptions::new();
        options.env_vars = self.env.clone();

        let script = self.script.compile(&ctx.variables()?)?;

        let args = vec![];
        let (code, output, _) = run(script.as_str(), &args, &options)?;

        if code != 0 {
            return Ok(RuleResult::Failure {
                name: SHELL_SCRIPT_NAME.into(),
                message: format!("exit code: {}", code),
            });
        }

        Ok(RuleResult::Success {
            name: SHELL_SCRIPT_NAME.into(),
            output: Some(output),
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::context::MockContext;
    use crate::rules::shell_script::ShellScriptRule;
    use crate::rules::{Rule, RuleResult};
    use crate::t;
    use anyhow::Result;
    use std::collections::HashMap;

    #[test]
    fn serialize_test() -> Result<()> {
        let config = ShellScriptRule {
            script: t!("echo hello"),
            env: None,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"script":"echo hello","env":null}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: ShellScriptRule = serde_json::from_str(r#"{"script":"echo hello"}"#)?;

        assert_eq!(config.script, t!("echo hello"));
        assert!(config.env.is_none());

        Ok(())
    }

    #[test]
    fn serialize_test_with_env() -> Result<()> {
        let mut env = HashMap::new();
        env.insert("TEST".to_string(), "value".to_string());

        let config = ShellScriptRule {
            script: t!("echo hello"),
            env: Some(env),
        };

        let serialized = serde_json::to_string(&config)?;

        assert!(serialized.contains("\"TEST\":\"value\""));
        assert!(serialized.contains("\"script\":\"echo hello\""));

        Ok(())
    }

    #[test]
    fn deserialize_test_with_env() -> Result<()> {
        let config: ShellScriptRule = serde_json::from_str(
            r#"{"script":"echo hello","env":{"TEST":"value"}}"#,
        )?;

        assert_eq!(config.script, t!("echo hello"));
        assert!(config.env.is_some());
        assert_eq!(config.env.unwrap().get("TEST"), Some(&"value".to_string()));

        Ok(())
    }




    fn mock_ctx_with_vars(vars: HashMap<String, String>) -> MockContext {
        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move || Ok(vars.clone()));
        ctx
    }

    fn mock_ctx() -> MockContext {
        mock_ctx_with_vars(HashMap::new())
    }

    #[test]
    fn test_shell_script() {
        let script = ShellScriptRule {
            script: t!("echo 'Test'"),
            env: None,
        };

        let result = script.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        #[cfg(not(windows))]
        assert_eq!(&output.unwrap(), "Test\n");
        #[cfg(windows)]
        assert_eq!(&output.unwrap(), "'Test'\r\n");
    }

    #[test]
    fn test_shell_script_failure() {
        let script = ShellScriptRule {
            script: t!("exit 1"),
            env: None,
        };

        let result = script.check(&mock_ctx()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "shell");
        assert_eq!(message, "exit code: 1");
    }

    #[test]
    fn test_shell_script_with_variables() {
        let mut variables = HashMap::new();
        variables.insert("name".to_string(), "Test".to_string());

        let script = ShellScriptRule {
            script: t!("echo 'Hello {{name}}'"),
            env: None,
        };

        let result = script.check(&mock_ctx_with_vars(variables)).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        #[cfg(not(windows))]
        assert_eq!(&output.unwrap(), "Hello Test\n");
        #[cfg(windows)]
        assert_eq!(&output.unwrap(), "'Hello Test'\r\n");
    }

    #[test]
    fn test_shell_script_with_env() {
        let mut env = HashMap::new();
        env.insert("TEST".to_string(), "Test".to_string());

        let script = ShellScriptRule {
            #[cfg(not(windows))]
            script: t!("echo $TEST"),
            #[cfg(windows)]
            script: t!("echo %TEST%"),
            env: Some(env),
        };

        let result = script.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        #[cfg(not(windows))]
        assert_eq!(&output.unwrap(), "Test\n");
        #[cfg(windows)]
        assert_eq!(&output.unwrap(), "Test\r\n");
    }

    #[test]
    fn test_shell_script_rule_new() {
        let rule = ShellScriptRule {
            script: "echo 'Test'".into(),
            env: None,
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        #[cfg(not(windows))]
        assert_eq!(&output.unwrap(), "Test\n");
        #[cfg(windows)]
        assert_eq!(&output.unwrap(), "'Test'\r\n");
    }

    #[test]
    fn test_shell_script_template_error() {
        let script = ShellScriptRule {
            env: None,
            script: t!("echo '{{missing}}'"),
        };

        let result = script.check(&mock_ctx());
        assert!(result.is_err());
    }

    #[test]
    fn test_display() {
        let rule = ShellScriptRule {
            script: t!("echo hello"),
            env: None,
        };
        assert_eq!(format!("{}", rule), "Run shell script: `echo hello`");
    }
}

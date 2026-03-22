use crate::context::Context;
use crate::rules::{Rule, RuleResult};
use crate::templates::{replace_in_hashmap, replace_in_vec};
use anyhow::Result;
use std::collections::HashMap;
use std::env;
use std::process::Command;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

#[derive(Debug, serde::Serialize, serde::Deserialize)]
pub struct ExecRule {
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
        let variables = ctx.variables_new()?;
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(
            &self.env.clone().unwrap_or_default(),
            &variables,
        )?);

        let output = Command::new(self.command.clone())
            .args(replace_in_vec(
                &self.args.clone().unwrap_or(vec![]),
                &variables,
            )?)
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
    use anyhow::Result;
    use assert2::assert;
    use std::collections::HashMap;

    #[cfg(windows)]
    static TEST_COMMAND: &str = "cmd";
    #[cfg(windows)]
    static TEST_COMMAN_ARGS: [&str; 3] = ["/C", "echo", "hello"];

    #[cfg(not(windows))]
    static TEST_COMMAND: &str = "echo";
    #[cfg(not(windows))]
    static TEST_COMMAN_ARGS: [&str; 1] = ["hello"];

    #[test]
    fn serialize_test() -> Result<()> {
        let config = ExecRule {
            command: "echo".to_string(),
            args: Some(vec!["hello".to_string()]),
            env: None,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"command":"echo","args":["hello"],"env":null}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test() -> Result<()> {
        let config: ExecRule = serde_json::from_str(
            r#"{"command":"echo","args":["hello"]}"#,
        )?;

        assert_eq!(config.command, "echo");
        assert_eq!(config.args, Some(vec!["hello".to_string()]));
        assert!(config.env.is_none());

        Ok(())
    }

    #[test]
    fn serialize_test_without_args() -> Result<()> {
        let config = ExecRule {
            command: "ls".to_string(),
            args: None,
            env: None,
        };

        let serialized = serde_json::to_string(&config)?;

        assert_eq!(
            serialized,
            r#"{"command":"ls","args":null,"env":null}"#
        );

        Ok(())
    }

    #[test]
    fn deserialize_test_without_args() -> Result<()> {
        let config: ExecRule = serde_json::from_str(r#"{"command":"ls"}"#)?;

        assert_eq!(config.command, "ls");
        assert!(config.args.is_none());
        assert!(config.env.is_none());

        Ok(())
    }

    #[test]
    fn serialize_test_with_env() -> Result<()> {
        let mut env = HashMap::new();
        env.insert("VAR".to_string(), "value".to_string());

        let config = ExecRule {
            command: "echo".to_string(),
            args: Some(vec!["hello".to_string()]),
            env: Some(env),
        };

        let serialized = serde_json::to_string(&config)?;

        assert!(serialized.contains("\"command\":\"echo\""));
        assert!(serialized.contains("\"VAR\":\"value\""));

        Ok(())
    }

    #[test]
    fn deserialize_test_with_env() -> Result<()> {
        let config: ExecRule = serde_json::from_str(
            r#"{"command":"echo","args":["hello"],"env":{"VAR":"value"}}"#,
        )?;

        assert_eq!(config.command, "echo");
        assert_eq!(config.args, Some(vec!["hello".to_string()]));
        assert!(config.env.is_some());
        assert_eq!(config.env.unwrap().get("VAR"), Some(&"value".to_string()));

        Ok(())
    }



    fn mock_ctx_with_vars(vars: HashMap<String, String>) -> MockContext {
        let mut ctx = MockContext::new();
        ctx.expect_variables_new().returning(move || Ok(vars.clone()));
        ctx
    }

    fn mock_ctx() -> MockContext {
        mock_ctx_with_vars(HashMap::new())
    }

    #[test]
    fn test_exec_rule() {
        let rule = ExecRule {
            command: TEST_COMMAND.into(),
            args: Some(test_args()),
            env: None,
        };

        let result = rule.check(&mock_ctx());

        let RuleResult::Success { name, output } = result.unwrap() else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        #[cfg(windows)]
        assert_eq!(output.unwrap(), "hello\r\n");
        #[cfg(not(windows))]
        assert_eq!(output.unwrap(), "hello\n");
    }

    #[test]
    fn test_exec_rule_with_env() {
        let mut env = HashMap::new();
        env.insert("HELLO".into(), "world".into());

        #[cfg(windows)]
        let (command, args) = ("cmd", vec!["/C", "echo %HELLO%"]);

        #[cfg(not(windows))]
        let (command, args) = ("sh", vec!["-c", "echo $HELLO"]);

        let rule = ExecRule {
            command: command.into(),
            args: Some(args.into_iter().map(String::from).collect()),
            env: Some(env),
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        #[cfg(not(windows))]
        assert_eq!(output.unwrap(), "world\n");
        #[cfg(windows)]
        assert_eq!(output.unwrap(), "world\r\n");
    }

    #[test]
    fn test_exec_rule_with_variables() {
        let mut vars = HashMap::new();
        vars.insert("HELLO".into(), "world".into());

        #[cfg(not(windows))]
        let rule = ExecRule {
            command: "echo".into(),
            args: Some(vec!["hello".into(), "{{HELLO}}".into()]),
            env: None,
        };

        #[cfg(windows)]
        let rule = ExecRule {
            command: "cmd".into(),
            args: Some(vec![
                "/C".into(),
                "echo".into(),
                "hello".into(),
                 "{{HELLO}}".into(),
            ]),
            env: None,
        };

        let result = rule.check(&mock_ctx_with_vars(vars)).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "exec");
        
        #[cfg(not(windows))]
        assert_eq!(output.unwrap(), "hello world\n");
        #[cfg(windows)]
        assert_eq!(output.unwrap(), "hello world\r\n");
    }

    #[test]
    fn test_exec_rule_failure_on_missing_file() {
        #[cfg(not(windows))]
        let rule = ExecRule {
            command: "cat".into(),
            args: Some(vec!["./unknown.txt".into()]),
            env: None,
        };

        #[cfg(windows)]
        let rule = ExecRule {
            command: "cmd".into(),
            args: Some(vec![
                "/C".into(),
                "type".into(),
                "unknown.txt".into(),
            ]),
            env: None,
        };

        let result = rule.check(&mock_ctx()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "exec");
        #[cfg(not(windows))]
        assert!(message.contains("No such file or directory"));
        #[cfg(windows)]
        assert!(message.contains("The system cannot find the file specified."));
    }

    #[test]
    fn test_return_error() {
        let rule = ExecRule {
            command: "XXXXXXXXXXXX".into(),
            args: None,
            env: None,
        };

        let result = rule.check(&mock_ctx()).err().unwrap();
        let error_msg = result.to_string();

        // Check for platform-specific error messages
        assert!(
            error_msg.contains("No such file or directory") || error_msg.contains("program not found"),
            "Error message should indicate command not found, got: {}",
            error_msg
        );
    }

    #[test]
    fn test_exec_rule_env_template_error() {
        let mut env = HashMap::new();
        env.insert("VAR".into(), "{{missing}}".into());

        let rule = ExecRule {
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
        };
        assert_eq!(format!("{}", rule), "Execute command: echo hello world");
    }

    #[test]
    fn test_display_without_args() {
        let rule = ExecRule {
            command: "ls".into(),
            args: None,
            env: None,
        };
        assert_eq!(format!("{}", rule), "Execute command: ls");
    }

    fn test_args() -> Vec<String> {
        return TEST_COMMAN_ARGS.iter()
            .map(|arg| arg.to_string())
            .collect()
    }
}

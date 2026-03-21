use crate::context::Context;
use crate::rules::rule::{ConditionalRule, Rule, RuleResult};
use crate::rules::{CompiledRule};
use crate::scripting::Expression;
use crate::templates::TemplateString;
use anyhow::Result;
use rules_derive::ConditionalRule as ConditionalRuleDerive;
use run_script::{run, ScriptOptions};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Serialize, Deserialize, ConditionalRuleDerive)]
pub struct ShellScriptRule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    pub script: TemplateString,
    pub env: Option<HashMap<String, String>>,
}

static SHELL_SCRIPT_NAME: &str = "shell";

#[typetag::serde(name = "shell")]
impl Rule for ShellScriptRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        if self.check_condition(ctx)? {
            return Ok(RuleResult::Skipped {
                name: SHELL_SCRIPT_NAME.into(),
            });
        }

        let mut options = ScriptOptions::new();
        options.env_vars = self.env.clone();

        let extract = self.extract.clone().unwrap_or(vec![]);
        let variables = ctx.variables(extract.as_slice())?;
        let script = self.script.compile(&variables)?;

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
    use crate::rules::rule::{Rule, RuleResult};
    use crate::rules::shell_script::{ShellScript, ShellScriptRule};
    use crate::rules::CompiledRule;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_shell_script() {
        let script = ShellScriptRule {
            when: None,
            extract: None,
            script: t!("echo 'Test'"),
            env: None,
        };

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_shell_script_failure() {
        let script = ShellScriptRule{
            when: None,
            extract: None,
            script: t!("exit 1"),
            env: None,
        };

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResult::Failure { name, message } = result else {
            unreachable!("Expected Failure");
        };
        assert_eq!(name, "Test");
        assert_eq!(message, "exit code: 1");
    }

    #[test]
    fn test_shell_script_with_variables() {
        let mut variables = HashMap::new();
        variables.insert("name".to_string(), "Test".to_string());

        let script = ShellScriptRule{
            when: None,
            extract: None,
            script: t!("echo 'Hello {{name}}'"),
            env: None,
        };

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Hello Test\n");
    }

    #[test]
    fn test_shell_script_with_env() {
        let mut env = HashMap::new();
        env.insert("TEST".to_string(), "Test".to_string());

        let script = ShellScriptRule {
            when: None,
            extract: None,
            script: t!("echo $TEST"),
            env: Some(env),
        };

        let result = script.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "Test");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_shell_script_rule_new() {
        let rule = ShellScriptRule {
            when: None,
            extract: None,
            script: "echo 'Test'".into(),
            env: None,
        };

        let result = rule.check(&MockContext::new()).unwrap();
        let RuleResult::Success { name, output } = result else {
            unreachable!("Expected Success");
        };
        assert_eq!(name, "shell");
        assert_eq!(output.unwrap(), "Test\n");
    }

    #[test]
    fn test_shell_script_template_error() {
        let script = ShellScriptRule {
            when: None,
            extract: None,
            env: None,
            script: t!("echo '{{missing}}'"),
        };

        let result = script.check(&MockContext::new());
        assert!(result.is_err());
    }
}

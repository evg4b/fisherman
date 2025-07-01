use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateString;
use anyhow::Result;
use run_script::{ScriptOptions, run};
use std::collections::HashMap;

pub struct ShellScript {
    name: String,
    script: TemplateString,
    env: HashMap<String, String>,
}

impl ShellScript {
    pub fn new(name: String, script: TemplateString, env: HashMap<String, String>) -> Self {
        Self { name, script, env }
    }
}

impl CompiledRule for ShellScript {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, ctx: &dyn Context) -> Result<RuleResult> {
        let mut options = ScriptOptions::new();
        options.env_vars = Some(self.env.clone());

        let args = vec![];
        let vars = ctx.variables(&[])?;
        let (code, output, _) = run(self.script.to_string(&vars)?.as_str(), &args, &options)?;

        if code != 0 {
            return Ok(RuleResult::Failure {
                name: self.name.clone(),
                message: format!("exit code: {}", code),
            });
        }

        Ok(RuleResult::Success {
            name: self.name.clone(),
            output,
        })
    }
}

#[cfg(test)]
mod tests {
    use crate::context::MockContext;
    use crate::rules::CompiledRule;
    use crate::rules::RuleResult;
    use crate::rules::shell_script::ShellScript;
    use crate::t;
    use std::collections::HashMap;

    #[test]
    fn test_shell_script() {
        let script = ShellScript::new("Test".to_string(), t!("echo 'Test'"), HashMap::new());

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        let result = script.check(&context).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(output, "Test\n");
    }

    #[test]
    fn test_shell_script_failure() {
        let script = ShellScript::new("Test".to_string(), t!("exit 1"), HashMap::new());

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        let result = script.check(&context).unwrap();
        let RuleResult::Failure { name, message } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(message, "exit code: 1");
    }

    #[test]
    fn test_shell_script_with_variables() {
        let script = ShellScript::new(
            "Test".to_string(),
            t!("echo 'Hello {{name}}'"),
            HashMap::new(),
        );

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| {
            let mut variables = HashMap::new();
            variables.insert("name".to_string(), "Test".to_string());
            Ok(variables)
        });

        let result = script.check(&context).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(output, "Hello Test\n");
    }

    #[test]
    fn test_shell_script_with_env() {
        let mut env = HashMap::new();
        env.insert("TEST".to_string(), "Test".to_string());

        let script = ShellScript::new("Test".to_string(), t!("echo $TEST"), env);

        let mut context = MockContext::new();
        context.expect_variables().returning(|_| Ok(HashMap::new()));

        let result = script.check(&context).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(output, "Test\n");
    }

    #[test]
    fn test_sync() {
        let script = ShellScript::new("Test".to_string(), t!("echo 'Test'"), HashMap::new());
        assert!(!script.sync());
    }
}

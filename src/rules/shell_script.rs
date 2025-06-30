use crate::context::Context;
use crate::rules::{CompiledRule, RuleResult};
use crate::templates::TemplateStringLegacy;
use anyhow::Result;
use run_script::{run, ScriptOptions};
use std::collections::HashMap;

pub struct ShellScript {
    name: String,
    script: TemplateStringLegacy,
    env: HashMap<String, String>,
}

impl ShellScript {
    pub fn new(name: String, script: TemplateStringLegacy, env: HashMap<String, String>) -> Self {
        Self { name, script, env }
    }
}

impl CompiledRule for ShellScript {
    fn sync(&self) -> bool {
        false
    }

    fn check(&self, _: &dyn Context) -> Result<RuleResult> {
        let mut options = ScriptOptions::new();
        options.env_vars = Some(self.env.clone());

        let args = vec![];
        let (code, output, _) = run(self.script.to_string()?.as_str(), &args, &options)?;

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
    use crate::rules::shell_script::ShellScript;
    use crate::rules::CompiledRule;
    use crate::rules::RuleResult;
    use crate::templates::TemplateStringLegacy;
    use std::collections::HashMap;

    #[test]
    fn test_shell_script() {
        let script = ShellScript::new(
            "Test".to_string(),
            TemplateStringLegacy::new("echo 'Test'".to_string(), HashMap::new()),
            HashMap::new(),
        );

        let context = MockContext::new();
        let result = script.check(&context).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(output, "Test\n");
    }

    #[test]
    fn test_shell_script_failure() {
        let script = ShellScript::new(
            "Test".to_string(),
            TemplateStringLegacy::new("exit 1".to_string(), HashMap::new()),
            HashMap::new(),
        );

        let context = MockContext::new();
        let result = script.check(&context).unwrap();
        let RuleResult::Failure { name, message } = result else {
            panic!("Rule failed")
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
            TemplateStringLegacy::new("echo 'Hello {{name}}'".to_string(), variables.clone()),
            HashMap::new(),
        );

        let context = MockContext::new();
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

        let script = ShellScript::new(
            "Test".to_string(),
            TemplateStringLegacy::new("echo $TEST".to_string(), HashMap::new()),
            env,
        );

        let context = MockContext::new();
        let result = script.check(&context).unwrap();
        let RuleResult::Success { name, output } = result else {
            panic!("Rule failed")
        };

        assert_eq!(name, "Test");
        assert_eq!(output, "Test\n");
    }
    
    #[test]
    fn test_sync() {
        let script = ShellScript::new(
            "Test".to_string(),
            TemplateStringLegacy::new("echo 'Test'".to_string(), HashMap::new()),
            HashMap::new(),
        );
        assert!(!script.sync());
    }
}

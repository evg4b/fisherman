use crate::rules::rule::{CompiledRule, RuleResult};
use crate::templates::replace_in_hashmap;
use anyhow::Result;
use std::collections::HashMap;
use std::env;
use std::process::Command;

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
    fn check(&self) -> Result<RuleResult> {
        let mut env_map: Env = env::vars().collect();
        env_map.extend(replace_in_hashmap(&self.env, &self.variables)?);

        let output = Command::new(self.command.clone())
            .args(self.args.clone())
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

use crate::common::BError;
use crate::rules::rule::{CompiledRule, RuleResult};
use crate::templates::replace_in_hashmap;
use std::collections::HashMap;
use std::env;
use std::process::Command;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

fn exec_rule(
    command: String,
    args: Args,
    env: Env,
    map: HashMap<String, String>,
) -> Result<(String, bool), BError> {
    let mut env_map: Env = env::vars().collect();
    env_map.extend(replace_in_hashmap(&env, &map)?);

    let cmd = Command::new(command).args(args).envs(env_map).output()?;

    Ok((
        String::from_utf8_lossy(&cmd.stdout).to_string(),
        cmd.status.success(),
    ))
}

#[derive(Debug)]
pub struct ExecRule {
    command: String,
    args: Args,
    env: Env,
    variables: HashMap<String, String>,
}

impl ExecRule {
    pub fn new(command: String, args: Args, env: Env, variables: HashMap<String, String>) -> Self {
        Self {
            command,
            args,
            env,
            variables,
        }
    }
}

impl CompiledRule for ExecRule {
    fn check(&self) -> RuleResult {
        match exec_rule(
            self.command.clone(),
            self.args.clone(),
            self.env.clone(),
            self.variables.clone(),
        ) {
            Ok((output, success)) => {
                if success {
                    RuleResult::Success {
                        name: self.command.clone(),
                        output,
                    }
                } else {
                    RuleResult::Failure {
                        name: self.command.clone(),
                        message: output,
                    }
                }
            }
            Err(e) => RuleResult::Failure {
                name: self.command.clone(),
                message: e.to_string(),
            },
        }
    }
}

use crate::common::BError;
use std::collections::HashMap;
use std::env;
use std::process::Command;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

pub fn exec_rule(command: String, args: Args, env: Env) -> Result<(String, bool), BError> {
    let mut env_map: Env = env::vars().collect();
    env_map.extend(env);

    let cmd = Command::new(command)
        .args(args)
        .envs(env_map)
        .output()?;

    Ok((
        String::from_utf8_lossy(&cmd.stdout).to_string(),
        cmd.status.success(),
    ))
}

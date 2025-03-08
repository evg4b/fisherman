use crate::common::BError;
use std::collections::HashMap;
use std::env;
use std::process::Command;

pub(crate) type Args = Vec<String>;
pub(crate) type Env = HashMap<String, String>;

pub fn exec_rule(command: String, args: Args, env: Env, map: HashMap<String, String>) -> Result<(String, bool), BError> {
    let mut env_map: Env = env::vars().collect();

    let transformed: HashMap<String, String> = env.iter()
        .map(|(k, v)| (k.clone(), replace_placeholders(v, &map)))  // Multiply each value by 10
        .collect();

    env_map.extend(transformed);

    let cmd = Command::new(command)
        .args(args)
        .envs(env_map)
        .output()?;

    Ok((
        String::from_utf8_lossy(&cmd.stdout).to_string(),
        cmd.status.success(),
    ))
}

fn replace_placeholders<T: Into<String>>(template: T, values: &HashMap<String, String>) -> String {
    let mut result: String = template.into();

    for (key, value) in values {
        let placeholder = format!("{{{{{}}}}}", key); // Creates "{{key}}"
        result = result.replace(&placeholder, value);
    }

    result
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_replace_placeholders() {
        let mut values = HashMap::new();
        values.insert("name".to_string(), "World".to_string());
        values.insert("greeting".to_string(), "Hello".to_string());

        let template = "{{greeting}}, {{name}}!";
        let result = replace_placeholders(template, &values);
        assert_eq!(result, "Hello, World!");
    }
}

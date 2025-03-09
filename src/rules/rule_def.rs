use crate::common::BError;
use crate::context::Context;
use crate::rules::exec_rule::ExecRule;
use crate::rules::rule::CompiledRule;
use crate::rules::vars::extract_variables;
use serde::Deserialize;
use std::collections::HashMap;

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub enum RuleRef {
    #[serde(rename = "exec")]
    ExecRule {
        command: String,
        args: Option<Vec<String>>,
        env: Option<HashMap<String, String>>,
        extract: Option<Vec<String>>,
    },
}

impl std::fmt::Display for RuleRef {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.name())
    }
}

impl RuleRef {
    pub(crate) fn name(&self) -> String {
        match self {
            RuleRef::ExecRule { command, args, .. } => {
                let args_str = args.as_ref().map_or(String::new(), |args| {
                    args.iter()
                        .map(|arg| {
                            if arg.contains(" ") {
                                format!("\"{}\"", arg.replace("\"", "\\\""))
                            } else {
                                arg.clone()
                            }
                        })
                        .collect::<Vec<String>>()
                        .join(" ")
                });
                format!("exec {} {}", command, args_str)
            }
        }
    }

    pub fn compile(&self, context: &impl Context, global_extract: Vec<String>) -> Result<impl CompiledRule, BError> {
        match self {
            RuleRef::ExecRule { command, args, env, extract } => {
                let mut local_extract = extract.clone().unwrap_or_default();
                local_extract.extend(global_extract);

                let variables = extract_variables(context, local_extract)?;
                Ok(ExecRule::new(
                    self.name(),
                    command.clone(),
                    args.clone().unwrap_or_default(),
                    env.clone().unwrap_or_default(),
                    variables,
                ))
            }
        }
    }
}

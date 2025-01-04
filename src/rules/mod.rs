pub mod exec_rule;
use crate::rules::exec_rule::exec_rule;
use serde::Deserialize;
use std::collections::HashMap;

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub(crate) enum RuleRef {
    #[serde(rename = "exec")]
    ExecRule {
        command: String,
        args: Option<Vec<String>>,
        env: Option<HashMap<String, String>>,
    },
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
                String::from(format!("exec {} {}", command, args_str))
            }
        }
    }
}

#[derive(Debug)]
pub(crate) struct Rule {
    def: RuleRef,
}

pub(crate) struct RuleResult {
    pub name: String,
    pub success: bool,
    pub message: String,
}

impl Rule {
    pub fn new(def: RuleRef) -> Self {
        Self { def }
    }

    pub fn exec(&self) -> RuleResult {
        match &self.def {
            RuleRef::ExecRule { command, args, env } => {
                match exec_rule(
                    command.clone(),
                    args.clone().unwrap_or_default(),
                    env.clone().unwrap_or_default(),
                ) {
                    Ok((message, success)) => RuleResult {
                        success,
                        message,
                        name: self.name(),
                    },
                    Err(e) => RuleResult {
                        success: false,
                        message: format!("Failed to execute rule: {}", e),
                        name: self.name(),
                    },
                }
            }
        }
    }

    fn name(&self) -> String {
        self.def.name()
    }
}

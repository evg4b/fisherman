use crate::context::Context;
use crate::rules::exec_rule::ExecRule;
use crate::rules::commit_message_regex::CommitMessageRegex;
use crate::rules::compiled_rule::CompiledRule;
use crate::rules::variables::extract_variables;
use anyhow::Result;
use serde::{Deserialize};
use std::collections::HashMap;
use crate::rules::commit_message_prefix::CommitMessagePrefix;
use crate::rules::commit_message_suffix::CommitMessageSuffix;
use crate::rules::shell_script::ShellScript;
use crate::scripting::Expression;

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub struct Rule {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    #[serde(flatten)]
    pub params: RuleParams,
}

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub enum RuleParams {
    #[serde(rename = "exec")]
    ExecRule {
        command: String,
        args: Option<Vec<String>>,
        env: Option<HashMap<String, String>>,
    },
    #[serde(rename = "message-regex")]
    CommitMessageRegex {
        regex: String,
    },
    #[serde(rename = "message-prefix")]
    CommitMessagePrefix {
        prefix: String,
    },
    #[serde(rename = "message-suffix")]
    CommitMessageSuffix {
        suffix: String,
    },
    #[serde(rename = "shell")]
    ShellScript {
        env: Option<HashMap<String, String>>,
        script: String,
    },
}

impl std::fmt::Display for Rule {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.params.name())
    }
}

macro_rules! wrap {
    ($expr:expr) => {
        Ok(Some(Box::new($expr)))
    };
}

impl Rule {
    pub fn compile(&self, context: &impl Context, global_extract: Vec<String>) -> Result<Option<Box<dyn CompiledRule>>> {
        let variables = prepare_variables(context, global_extract, &self.extract)?;

        if let Some(expression) = &self.when {
            if !expression.check(&variables)? {
                return Ok(None);
            }
        }

        match &self.params {
            RuleParams::ExecRule { command, args, env,  .. } => {
                wrap!(ExecRule::new(
                    self.to_string(),
                    command.clone(),
                    args.clone().unwrap_or_default(),
                    env.clone().unwrap_or_default(),
                    variables,
                ))
            }
            RuleParams::CommitMessageRegex { regex, .. } => {
                wrap!(CommitMessageRegex::new(
                    self.to_string(),
                    regex.clone(),
                    variables,
                ))
            },
            RuleParams::CommitMessagePrefix { prefix,  .. } => {
                wrap!(CommitMessagePrefix::new(
                    self.to_string(),
                    prefix.clone(),
                    variables,
                ))
            },
            RuleParams::CommitMessageSuffix { suffix,.. } => {
                wrap!(CommitMessageSuffix::new(
                    self.to_string(),
                    suffix.clone(),
                    variables,
                ))
            },
            RuleParams::ShellScript { script, env, .. } => {
                wrap!(ShellScript::new(
                    self.to_string(),
                    script.clone(),
                    env.clone().unwrap_or_default(),
                    variables,
                ))
            }
        }
    }
}

impl RuleParams {
    pub(crate) fn name(&self) -> String {
        match self {
            RuleParams::ExecRule { command, args, .. } => {
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
            RuleParams::CommitMessageRegex { regex, .. } => {
                format!("commit message rule should match regex: {}", regex)
            }
            RuleParams::CommitMessagePrefix { prefix, .. } => {
                format!("commit message rule should start with: {}", prefix)
            },
            RuleParams::CommitMessageSuffix { suffix, .. } => {
                format!("commit message rule should end with: {}", suffix)
            },
            RuleParams::ShellScript { script, .. } => {
                format!("shell script:\n{}", script)
            }
        }
    }
}

fn prepare_variables(context: &impl Context, global: Vec<String>, local: &Option<Vec<String>>) -> Result<HashMap<String, String>> {
    let mut variables = local.clone().unwrap_or_default();
    variables.extend(global);
    extract_variables(context, variables)
}
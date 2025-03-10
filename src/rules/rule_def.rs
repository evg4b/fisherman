use crate::context::Context;
use crate::rules::exec_rule::ExecRule;
use crate::rules::commit_message_regex::CommitMessageRegex;
use crate::rules::compiled_rule::CompiledRule;
use crate::rules::variables::extract_variables;
use anyhow::Result;
use serde::Deserialize;
use std::collections::HashMap;
use regex::Regex;
use crate::rules::commit_message_prefix::CommitMessagePrefix;
use crate::rules::commit_message_suffix::CommitMessageSuffix;

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub enum RuleDefinition {
    #[serde(rename = "exec")]
    ExecRule {
        command: String,
        args: Option<Vec<String>>,
        env: Option<HashMap<String, String>>,
        extract: Option<Vec<String>>,
    },
    #[serde(rename = "message-regex")]
    CommitMessageRegex {
        #[serde(rename = "regex")]
        regex: String,
    },
    #[serde(rename = "message-prefix")]
    CommitMessagePrefix {
        extract: Option<Vec<String>>,
        #[serde(rename = "prefix")]
        prefix: String,
    },
    #[serde(rename = "message-suffix")]
    CommitMessageSuffix {
        extract: Option<Vec<String>>,
        #[serde(rename = "suffix")]
        suffix: String,
    },
}

impl std::fmt::Display for RuleDefinition {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.name())
    }
}

impl RuleDefinition {
    pub(crate) fn name(&self) -> String {
        match self {
            RuleDefinition::ExecRule { command, args, .. } => {
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
            RuleDefinition::CommitMessageRegex { regex } => {
                format!("commit message rule should match regex: {}", regex)
            }
            RuleDefinition::CommitMessagePrefix { prefix, .. } => {
                format!("commit message rule should start with: {}", prefix)
            },
            RuleDefinition::CommitMessageSuffix { suffix, .. } => {
                format!("commit message rule should end with: {}", suffix)
            }
        }
    }

    pub fn compile(&self, context: &impl Context, global_extract: Vec<String>) -> Result<Box<dyn CompiledRule>> {
        match self {
            RuleDefinition::ExecRule { command, args, env, extract } => {
                let rule = ExecRule::new(
                    self.name(),
                    command.clone(),
                    args.clone().unwrap_or_default(),
                    env.clone().unwrap_or_default(),
                    prepare_variables(context, global_extract, extract)?,
                );

                Ok(Box::new(rule))
            }
            RuleDefinition::CommitMessageRegex { regex } => {
                let rule = CommitMessageRegex::new(
                    self.name(),
                    Regex::new(regex)?,
                );

                Ok(Box::new(rule))
            },
            RuleDefinition::CommitMessagePrefix { prefix, extract } => {
                let rule = CommitMessagePrefix::new(
                    self.name(),
                    prefix.clone(),
                    prepare_variables(context, global_extract, extract)?,
                );

                Ok(Box::new(rule))
            },
            RuleDefinition::CommitMessageSuffix { suffix, extract } => {
                let rule = CommitMessageSuffix::new(
                    self.name(),
                    suffix.clone(),
                    prepare_variables(context, global_extract, extract)?,
                );

                Ok(Box::new(rule))
            }
        }
    }
}

fn prepare_variables(context: &impl Context, global: Vec<String>, local: &Option<Vec<String>>) -> Result<HashMap<String, String>> {
    let mut variables = local.clone().unwrap_or_default();
    variables.extend(global);
    extract_variables(context, variables)
}
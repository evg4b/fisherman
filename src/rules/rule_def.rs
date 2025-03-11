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
use crate::rules::shell_script::ShellScript;

#[derive(Debug, Deserialize, Clone)]
#[serde(tag = "type")]
pub struct Rule {
    pub when: Option<String>,
    #[serde(flatten)]
    pub rule: RuleDefinition,
}

#[cfg(test)]
mod rule_test {
    #[test]
    fn test_rule() {
        let rule = r#"
        {
            "when": "is_def_var(\"xx\") && parse_int(xx) > 10",
            "type": "exec",
            "command": "echo",
            "args": ["hello", "world"],
            "env": {
                "xx": "20"
            },
            "extract": ["yy"]
        }
        "#;

        let rule: super::Rule = serde_json::from_str(rule).unwrap();
        assert_eq!(rule.when, Some("is_def_var(\"xx\") && parse_int(xx) > 10".to_string()));
    }
}


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
        regex: String,
    },
    #[serde(rename = "message-prefix")]
    CommitMessagePrefix {
        extract: Option<Vec<String>>,
        prefix: String,
    },
    #[serde(rename = "message-suffix")]
    CommitMessageSuffix {
        extract: Option<Vec<String>>,
        suffix: String,
    },
    #[serde(rename = "shell")]
    ShellScript {
        extract: Option<Vec<String>>,
        env: Option<HashMap<String, String>>,
        script: String,
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
            RuleDefinition::CommitMessageRegex { regex, .. } => {
                format!("commit message rule should match regex: {}", regex)
            }
            RuleDefinition::CommitMessagePrefix { prefix, .. } => {
                format!("commit message rule should start with: {}", prefix)
            },
            RuleDefinition::CommitMessageSuffix { suffix, .. } => {
                format!("commit message rule should end with: {}", suffix)
            },
            RuleDefinition::ShellScript { script, .. } => {
                format!("shell script:\n{}", script)
            }
        }
    }

    pub fn compile(&self, context: &impl Context, global_extract: Vec<String>) -> Result<Box<dyn CompiledRule>> {
        match self {
            RuleDefinition::ExecRule { command, args, env, extract, .. } => {
                Ok(Box::new(ExecRule::new(
                    self.name(),
                    command.clone(),
                    args.clone().unwrap_or_default(),
                    env.clone().unwrap_or_default(),
                    prepare_variables(context, global_extract, extract)?,
                )))
            }
            RuleDefinition::CommitMessageRegex { regex, .. } => {
                Ok(Box::new(CommitMessageRegex::new(
                    self.name(),
                    Regex::new(regex)?,
                )))
            },
            RuleDefinition::CommitMessagePrefix { prefix, extract, .. } => {
                Ok(Box::new(CommitMessagePrefix::new(
                    self.name(),
                    prefix.clone(),
                    prepare_variables(context, global_extract, extract)?,
                )))
            },
            RuleDefinition::CommitMessageSuffix { suffix, extract,.. } => {
                Ok(Box::new(CommitMessageSuffix::new(
                    self.name(),
                    suffix.clone(),
                    prepare_variables(context, global_extract, extract)?,
                )))
            },
            RuleDefinition::ShellScript { script, extract, env, .. } => {
                Ok(Box::new(ShellScript::new(
                    self.name(),
                    script.clone(),
                    env.clone().unwrap_or_default(),
                    prepare_variables(context, global_extract, extract)?,
                )))
            }
        }
    }
}

fn prepare_variables(context: &impl Context, global: Vec<String>, local: &Option<Vec<String>>) -> Result<HashMap<String, String>> {
    let mut variables = local.clone().unwrap_or_default();
    variables.extend(global);
    extract_variables(context, variables)
}
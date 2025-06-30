use crate::context::Context;
use crate::rules::branch_name_prefix::BranchNamePrefix;
use crate::rules::branch_name_regex::BranchNameRegex;
use crate::rules::branch_name_suffix::BranchNameSuffix;
use crate::rules::commit_message_prefix::CommitMessagePrefix;
use crate::rules::commit_message_regex::CommitMessageRegex;
use crate::rules::commit_message_suffix::CommitMessageSuffix;
use crate::rules::compiled_rule::CompiledRule;
use crate::rules::exec_rule::ExecRule;
use crate::rules::shell_script::ShellScript;
use crate::rules::write_file::WriteFile;
use crate::scripting::Expression;
use crate::templates::TemplateString;
use crate::tmpl_legacy;
use anyhow::Result;
use serde::Deserialize;
use std::collections::HashMap;

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
    CommitMessageRegex { regex: String },
    #[serde(rename = "message-prefix")]
    CommitMessagePrefix { prefix: String },
    #[serde(rename = "message-suffix")]
    CommitMessageSuffix { suffix: String },
    #[serde(rename = "shell")]
    ShellScript {
        env: Option<HashMap<String, String>>,
        script: String,
    },
    #[serde(rename = "write-file")]
    WriteFile {
        path: String,
        content: String,
        append: Option<bool>,
    },
    #[serde(rename = "branch-name-regex")]
    BranchNameRegex { regex: String },
    #[serde(rename = "branch-name-prefix")]
    BranchNamePrefix { prefix: String },
    #[serde(rename = "branch-name-suffix")]
    BranchNameSuffix { suffix: String },
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
    pub fn compile(&self, context: &impl Context) -> Result<Option<Box<dyn CompiledRule>>> {
        let variables = context.variables(self.extract.as_ref().unwrap_or(&vec![]))?;

        if let Some(expression) = &self.when {
            if !expression.check(&variables)? {
                return Ok(None);
            }
        }

        match &self.params {
            RuleParams::ExecRule {
                command, args, env, ..
            } => {
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
                    tmpl_legacy!(regex.clone(), variables)
                ))
            }
            RuleParams::CommitMessagePrefix { prefix, .. } => {
                wrap!(CommitMessagePrefix::new(
                    self.to_string(),
                    tmpl_legacy!(prefix.clone(), variables),
                ))
            }
            RuleParams::CommitMessageSuffix { suffix, .. } => {
                wrap!(CommitMessageSuffix::new(
                    self.to_string(),
                    tmpl_legacy!(suffix.clone(), variables),
                ))
            }
            RuleParams::ShellScript { script, env, .. } => {
                wrap!(ShellScript::new(
                    self.to_string(),
                    tmpl_legacy!(script.clone(), variables.clone()),
                    env.clone().unwrap_or_default(),
                ))
            }
            RuleParams::WriteFile {
                path,
                content,
                append,
            } => {
                wrap!(WriteFile::new(
                    self.to_string(),
                    TemplateString::from(path.clone()),
                    TemplateString::from(content.clone()),
                    append.unwrap_or(false),
                ))
            }
            RuleParams::BranchNameRegex { regex, .. } => {
                wrap!(BranchNameRegex::new(
                    self.to_string(),
                    tmpl_legacy!(regex.clone(), variables)
                ))
            }
            RuleParams::BranchNamePrefix { prefix, .. } => {
                wrap!(BranchNamePrefix::new(
                    self.to_string(),
                    tmpl_legacy!(prefix.clone(), variables),
                ))
            }
            RuleParams::BranchNameSuffix { suffix, .. } => {
                wrap!(BranchNameSuffix::new(
                    self.to_string(),
                    tmpl_legacy!(suffix.clone(), variables),
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
            }
            RuleParams::CommitMessageSuffix { suffix, .. } => {
                format!("commit message rule should end with: {}", suffix)
            }
            RuleParams::ShellScript { script, .. } => {
                format!("shell script:\n{}", script)
            }
            RuleParams::WriteFile { path, .. } => {
                format!("write file to: {}", path)
            }
            RuleParams::BranchNameRegex { regex, .. } => {
                format!("branch name rule should match regex: {}", regex)
            }
            RuleParams::BranchNamePrefix { prefix, .. } => {
                format!("branch name rule should start with: {}", prefix)
            }
            RuleParams::BranchNameSuffix { suffix, .. } => {
                format!("branch name rule should end with: {}", suffix)
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_exec_rule_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::ExecRule {
                command: "echo".to_string(),
                args: Some(vec!["Hello".to_string(), "World".to_string()]),
                env: None,
            },
        };
        assert_eq!(rule.to_string(), "exec echo Hello World");
    }

    #[test]
    fn test_exec_rule_with_params_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::ExecRule {
                command: "echo".to_string(),
                args: Some(vec!["Hello".to_string(), "World".to_string()]),
                env: Some(HashMap::from([("KEY".to_string(), "VALUE".to_string())])),
            },
        };
        assert_eq!(rule.to_string(), "exec echo Hello World");
    }

    #[test]
    fn test_commit_message_regex_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessageRegex {
                regex: r"^feat\(".to_string(),
            },
        };
        assert_eq!(
            rule.to_string(),
            "commit message rule should match regex: ^feat\\("
        );
    }

    #[test]
    fn test_commit_message_prefix_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat: ".to_string(),
            },
        };
        assert_eq!(
            rule.to_string(),
            "commit message rule should start with: feat: "
        );
    }

    #[test]
    fn test_commit_message_suffix_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessageSuffix {
                suffix: "Fixes #123".to_string(),
            },
        };
        assert_eq!(
            rule.to_string(),
            "commit message rule should end with: Fixes #123"
        );
    }

    #[test]
    fn test_shell_script_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::ShellScript {
                script: "echo 'Hello, World!'".to_string(),
                env: None,
            },
        };
        assert_eq!(rule.to_string(), "shell script:\necho 'Hello, World!'");
    }

    #[test]
    fn test_write_file_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::WriteFile {
                path: "/tmp/test.txt".to_string(),
                content: "Hello, World!".to_string(),
                append: Some(false),
            },
        };
        assert_eq!(rule.to_string(), "write file to: /tmp/test.txt");
    }

    #[test]
    fn test_branch_name_regex_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNameRegex {
                regex: r"^feature/".to_string(),
            },
        };
        assert_eq!(
            rule.to_string(),
            "branch name rule should match regex: ^feature/"
        );
    }

    #[test]
    fn test_branch_name_prefix_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNamePrefix {
                prefix: "feature/".to_string(),
            },
        };
        assert_eq!(
            rule.to_string(),
            "branch name rule should start with: feature/"
        );
    }

    #[test]
    fn test_branch_name_suffix_display() {
        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNameSuffix {
                suffix: "/v1".to_string(),
            },
        };
        assert_eq!(rule.to_string(), "branch name rule should end with: /v1");
    }
}

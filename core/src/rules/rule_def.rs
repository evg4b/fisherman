use crate::context::Context;
use crate::rules::compiled_rule::CompiledRule;
use crate::scripting::Expression;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct RuleOLD {
    pub when: Option<Expression>,
    pub extract: Option<Vec<String>>,
    #[serde(flatten)]
    pub params: RuleParams,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
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
    #[serde(rename = "copy-files")]
    CopyFiles {
        glob: String,
        destination: String,
        source: Option<String>,
    },
    #[serde(rename = "delete-files")]
    DeleteFiles {
        glob: String,
        #[serde(rename = "fail-if-not-found")]
        fail_if_not_found: Option<bool>,
    },
    #[serde(rename = "suppress-files")]
    SuppressFiles { glob: String },
    #[serde(rename = "suppress-string")]
    SuppressString { regex: String, glob: Option<String> },
}

impl std::fmt::Display for RuleOLD {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.params.name())
    }
}

macro_rules! wrap {
    ($expr:expr) => {
        Ok(Some(Box::new($expr)))
    };
}

impl RuleParams {
    pub(crate) fn name(&self) -> String {
        match self {
            RuleParams::ExecRule { command, args, .. } => {
                let args_str = args.as_ref().map_or(String::new(), |args| {
                    args.iter()
                        .map(|arg| {
                            if arg.contains(' ') {
                                format!("\"{}\"", arg.replace('"', "\\\""))
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
                format!("write a file to: {}", path)
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
            RuleParams::CopyFiles {
                glob,
                destination,
                source,
            } => {
                format!(
                    "copy files from {} to {} (source: {})",
                    glob,
                    destination,
                    source.as_ref().unwrap_or(&String::from("<n/a>"))
                )
            }
            RuleParams::DeleteFiles {
                glob,
                fail_if_not_found,
            } => {
                format!(
                    "delete files matching {} {}",
                    glob,
                    fail_if_not_found.unwrap_or(false)
                )
            }
            RuleParams::SuppressFiles { glob } => {
                format!("suppress files matching {}", glob)
            }
            RuleParams::SuppressString { regex, glob } => {
                format!(
                    "suppress string matching {} in {}",
                    regex,
                    glob.as_deref().unwrap_or("*")
                )
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_exec_rule_display() {
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
            when: None,
            extract: None,
            params: RuleParams::WriteFile {
                path: "/tmp/test.txt".to_string(),
                content: "Hello, World!".to_string(),
                append: Some(false),
            },
        };
        assert_eq!(rule.to_string(), "write a file to: /tmp/test.txt");
    }

    #[test]
    fn test_branch_name_regex_display() {
        let rule = RuleOLD {
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
        let rule = RuleOLD {
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
        let rule = RuleOLD {
            when: None,
            extract: None,
            params: RuleParams::BranchNameSuffix {
                suffix: "/v1".to_string(),
            },
        };
        assert_eq!(rule.to_string(), "branch name rule should end with: /v1");
    }

    #[test]
    fn test_rule_params_name_exec_with_args_containing_spaces() {
        let params = RuleParams::ExecRule {
            command: "echo".to_string(),
            args: Some(vec!["Hello World".to_string(), "test".to_string()]),
            env: None,
        };
        assert_eq!(params.name(), "exec echo \"Hello World\" test");
    }

    #[test]
    fn test_rule_params_name_exec_with_args_containing_quotes() {
        let params = RuleParams::ExecRule {
            command: "echo".to_string(),
            args: Some(vec!["Say \"Hello\"".to_string()]),
            env: None,
        };
        assert_eq!(params.name(), "exec echo \"Say \\\"Hello\\\"\"");
    }

    #[test]
    fn test_rule_params_name_exec_without_args() {
        let params = RuleParams::ExecRule {
            command: "echo".to_string(),
            args: None,
            env: None,
        };
        assert_eq!(params.name(), "exec echo ");
    }

    #[test]
    fn test_rule_params_name_copy_files_with_source() {
        let params = RuleParams::CopyFiles {
            glob: "*.txt".to_string(),
            destination: "dest".to_string(),
            source: Some("src".to_string()),
        };
        assert_eq!(params.name(), "copy files from *.txt to dest (source: src)");
    }

    #[test]
    fn test_rule_params_name_copy_files_without_source() {
        let params = RuleParams::CopyFiles {
            glob: "*.txt".to_string(),
            destination: "dest".to_string(),
            source: None,
        };
        assert_eq!(
            params.name(),
            "copy files from *.txt to dest (source: <n/a>)"
        );
    }

    #[test]
    fn test_rule_params_name_delete_files() {
        let params = RuleParams::DeleteFiles {
            glob: "*.log".to_string(),
            fail_if_not_found: Some(true),
        };
        assert_eq!(params.name(), "delete files matching *.log true");
    }
}

use crate::context::Context;
use crate::rules::branch_name_prefix::BranchNamePrefix;
use crate::rules::branch_name_regex::BranchNameRegex;
use crate::rules::branch_name_suffix::BranchNameSuffix;
use crate::rules::commit_message_prefix::CommitMessagePrefix;
use crate::rules::commit_message_regex::CommitMessageRegex;
use crate::rules::commit_message_suffix::CommitMessageSuffix;
use crate::rules::compiled_rule::CompiledRule;
use crate::rules::copy_files::CopyFiles;
use crate::rules::delete_files::DeleteFiles;
use crate::rules::exec_rule::ExecRule;
use crate::rules::shell_script::ShellScript;
use crate::rules::write_file::WriteFile;
use crate::scripting::Expression;
use crate::t;
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
    #[serde(rename = "write-files")]
    CopyFiles { glob: String, destination: String },
    #[serde(rename = "delete-files")]
    DeleteFiles {
        glob: String,
        #[serde(rename = "fail-if-not-found")]
        fail_if_not_found: Option<bool>,
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
    pub fn compile(&self, context: &impl Context) -> Result<Option<Box<dyn CompiledRule>>> {
        let variables = context.variables(self.extract.as_ref().unwrap_or(&vec![]))?;

        if let Some(expression) = &self.when && !expression.check(&variables)? {
            return Ok(None);
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
                    t!(regex.clone())
                ))
            }
            RuleParams::CommitMessagePrefix { prefix, .. } => {
                wrap!(CommitMessagePrefix::new(
                    self.to_string(),
                    t!(prefix.clone()),
                ))
            }
            RuleParams::CommitMessageSuffix { suffix, .. } => {
                wrap!(CommitMessageSuffix::new(
                    self.to_string(),
                    t!(suffix.clone()),
                ))
            }
            RuleParams::ShellScript { script, env, .. } => {
                wrap!(ShellScript::new(
                    self.to_string(),
                    t!(script.clone()),
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
                    t!(path.clone()),
                    t!(content.clone()),
                    append.unwrap_or(false),
                ))
            }
            RuleParams::BranchNameRegex { regex, .. } => {
                wrap!(BranchNameRegex::new(self.to_string(), t!(regex.clone())))
            }
            RuleParams::BranchNamePrefix { prefix, .. } => {
                wrap!(BranchNamePrefix::new(
                    self.to_string(),
                    t!(prefix.clone()),
                ))
            }
            RuleParams::BranchNameSuffix { suffix, .. } => {
                wrap!(BranchNameSuffix::new(
                    self.to_string(),
                    t!(suffix.clone()),
                ))
            }
            RuleParams::CopyFiles { glob, destination } => {
                wrap!(CopyFiles::new(
                    self.to_string(),
                    t!(glob.clone()),
                    t!(destination.clone()),
                ))
            }
            RuleParams::DeleteFiles { glob, fail_if_not_found } => {
                wrap!(DeleteFiles::new(
                    self.to_string(),
                    t!(glob.clone()),
                    fail_if_not_found.unwrap_or(false),
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
            RuleParams::CopyFiles { glob, destination } => {
                format!("copy files from {} to {}", glob, destination)
            }
            RuleParams::DeleteFiles { glob, fail_if_not_found } => {
                format!("delete files matching {} {}", glob, fail_if_not_found.unwrap_or(false))
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
        assert_eq!(rule.to_string(), "write a file to: /tmp/test.txt");
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

    #[test]
    fn test_compile_exec_rule() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::ExecRule {
                command: "echo".to_string(),
                args: Some(vec!["test".to_string()]),
                env: Some(HashMap::from([("KEY".to_string(), "VALUE".to_string())])),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_with_when_condition_true() {
        use crate::context::MockContext;
        use crate::scripting::Expression;

        let mut vars = HashMap::new();
        vars.insert("branch".to_string(), "main".to_string());

        let rule = Rule {
            when: Some(Expression::new("branch == \"main\"")),
            extract: None,
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_with_when_condition_false() {
        use crate::context::MockContext;
        use crate::scripting::Expression;

        let mut vars = HashMap::new();
        vars.insert("branch".to_string(), "develop".to_string());

        let rule = Rule {
            when: Some(Expression::new("branch == \"main\"")),
            extract: None,
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables().returning(move |_| Ok(vars.clone()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_none());
    }

    #[test]
    fn test_compile_with_extract() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: Some(vec!["branch".to_string(), "ticket".to_string()]),
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_commit_message_regex() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessageRegex {
                regex: "^feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_commit_message_suffix() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessageSuffix {
                suffix: "done".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_shell_script() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::ShellScript {
                script: "echo test".to_string(),
                env: Some(HashMap::from([("KEY".to_string(), "VALUE".to_string())])),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_write_file() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::WriteFile {
                path: "/tmp/test.txt".to_string(),
                content: "content".to_string(),
                append: Some(true),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_branch_name_regex() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNameRegex {
                regex: "^feature/".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_branch_name_prefix() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNamePrefix {
                prefix: "feature/".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
    }

    #[test]
    fn test_compile_branch_name_suffix() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::BranchNameSuffix {
                suffix: "/v1".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx).unwrap();
        assert!(result.is_some());
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
    fn test_compile_variables_error() {
        use crate::context::MockContext;

        let rule = Rule {
            when: None,
            extract: None,
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Err(anyhow::anyhow!("Variables error")));

        let result = rule.compile(&ctx);
        assert!(result.is_err());
    }

    #[test]
    fn test_compile_when_expression_error() {
        use crate::context::MockContext;
        use crate::scripting::Expression;

        let rule = Rule {
            when: Some(Expression::new("invalid expression !!!")),
            extract: None,
            params: RuleParams::CommitMessagePrefix {
                prefix: "feat".to_string(),
            },
        };

        let mut ctx = MockContext::new();
        ctx.expect_variables()
            .returning(|_| Ok(HashMap::<String, String>::new()));

        let result = rule.compile(&ctx);
        assert!(result.is_err());
    }
}

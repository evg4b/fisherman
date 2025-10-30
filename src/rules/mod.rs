mod branch_name_prefix;
mod branch_name_regex;
mod branch_name_suffix;
mod commit_message_prefix;
mod commit_message_regex;
mod commit_message_suffix;
mod compiled_rule;
mod exec_rule;
mod helpers;
mod rule_def;
mod shell_script;
mod write_file;
mod copy_files;
mod delete_files;

pub use crate::rules::compiled_rule::CompiledRule;
pub use crate::rules::compiled_rule::RuleResult;
pub use crate::rules::rule_def::Rule;
#[allow(unused_imports)]
pub use crate::rules::rule_def::RuleParams;

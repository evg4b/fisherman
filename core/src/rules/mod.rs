pub mod branch_name_prefix;
pub mod branch_name_regex;
pub mod branch_name_suffix;
pub mod commit_message_prefix;
pub mod commit_message_regex;
pub mod commit_message_suffix;
mod compiled_rule;
mod copy_files;
mod delete_files;
pub mod exec_rule;
mod helpers;
pub mod rule;
mod rule_def;
pub mod shell_script;
mod suppress_files;
mod suppress_string;
pub mod write_file;

pub use crate::rules::compiled_rule::{CompiledRule, RuleResultOld};
#[allow(unused_imports)]
pub use crate::rules::rule_def::{RuleOLD, RuleParams};

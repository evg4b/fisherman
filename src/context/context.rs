use std::collections::HashMap;
use crate::configuration::Configuration;
use anyhow::Result;
use mockall::automock;
use std::path::{Path, PathBuf};

#[automock]
pub trait Context {
    fn repo_path(&self) -> &Path;
    fn hooks_dir(&self) -> PathBuf;
    fn bin(&self) -> &Path;
    fn current_branch(&self) -> Result<String>;
    fn commit_msg(&self) -> Result<String>;
    fn set_commit_msg_path(&mut self, message_file: PathBuf);
    fn configuration(&self) -> Result<Configuration>;
    fn variables(&self, additional: &Vec<String>) ->  Result<HashMap<String, String>>;
}

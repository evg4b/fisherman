use crate::configuration::Configuration;
use anyhow::Result;
use mockall::automock;
use std::collections::HashMap;
use std::path::{Path, PathBuf};

#[automock]
pub trait Context: Send {
    fn repo_path(&self) -> &Path;
    fn hooks_dir(&self) -> PathBuf;
    fn bin(&self) -> &Path;
    fn current_branch(&self) -> Result<String>;
    fn commit_msg(&self) -> Result<String>;
    fn set_commit_msg_path(&mut self, message_file: PathBuf);
    fn configuration(&self) -> Result<Configuration>;
    fn variables(&self, additional: &[String]) -> Result<HashMap<String, String>>;
}

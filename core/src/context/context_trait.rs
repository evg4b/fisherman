use crate::configuration::Configuration;
use anyhow::Result;
use mockall::automock;
use std::collections::HashMap;
use std::path::{Path, PathBuf};

#[derive(Debug, Clone, PartialEq)]
pub enum DiffLine {
    Added(String),
    Deleted(String),
}

#[automock]
pub trait Context: Send + Sync {
    fn repo_path(&self) -> &Path;
    fn hooks_dir(&self) -> PathBuf;
    fn bin(&self) -> &Path;
    fn current_branch(&self) -> Result<String>;
    fn commit_msg(&self) -> Result<String>;
    fn set_commit_msg_path(&mut self, message_file: PathBuf);
    fn configuration(&self) -> Result<Configuration>;
    fn variables(&self, extract: &[String]) -> Result<HashMap<String, String>>;
    fn staged_files(&self) -> Result<Vec<PathBuf>>;
    fn staged_diff(&self, path: &Path) -> Result<Vec<DiffLine>>;
}

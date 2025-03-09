use std::path::{Path, PathBuf};
use mockall::automock;
use crate::common::BError;

#[automock]
pub trait Context {
    fn repo_path(&self) -> &Path;
    fn hooks_dir(&self) -> PathBuf;
    fn bin(&self) -> &Path;
    fn current_branch(&self) -> Result<String, BError>;
}
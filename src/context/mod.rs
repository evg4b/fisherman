use std::fmt::Display;
use crate::common::BError;
use git2::Repository;
use std::path::{Path, PathBuf};

pub(crate) struct Context {
    repo: Repository,
    cwd: PathBuf,
    bin: PathBuf,
}

impl Context {
    pub(crate) fn new(cwd: PathBuf) -> Result<Self, BError> {
        let repo = Repository::open(cwd.clone())?;
        let bin = std::env::current_exe()?;

        Ok(Self { repo, cwd, bin })
    }

    pub(crate) fn current_branch(&self) -> Result<String, BError> {
        let head = self.repo.head()?;
        Ok(head.shorthand().unwrap_or("HEAD").to_string())
    }
    
    pub(crate) fn repo_path(&self) -> &Path {
        self.cwd.as_path()
    }

    pub(crate) fn hooks_dir(&self) -> PathBuf {
        self.cwd.join(".git/hooks")
    }

    pub(crate) fn bin(&self) -> &Path {
        self.bin.as_path()
    }
}

impl Display for Context {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Context {{ repo: {:?}, cwd: {:?} }}", self.repo.path(), self.cwd)
    }
}
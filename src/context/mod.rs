use crate::common::BError;
use git2::Repository;
use std::fmt::Display;
use std::path::{Path, PathBuf};

pub(crate) struct GitRepoContext {
    repo: Repository,
    cwd: PathBuf,
    bin: PathBuf,
}

pub(crate) trait Context {
    fn repo_path(&self) -> &Path;
    fn hooks_dir(&self) -> PathBuf;
    fn bin(&self) -> &Path;
    fn current_branch(&self) -> Result<String, BError>;
}

impl Context for GitRepoContext {
    fn repo_path(&self) -> &Path {
        self.cwd.as_path()
    }

    fn hooks_dir(&self) -> PathBuf {
        self.cwd.join(".git/hooks")
    }

    fn bin(&self) -> &Path {
        self.bin.as_path()
    }

    fn current_branch(&self) -> Result<String, BError> {
        let head = self.repo.head()?;
        Ok(head.shorthand().unwrap_or("HEAD").to_string())
    }
}

impl GitRepoContext {
    pub(crate) fn new(cwd: PathBuf) -> Result<Self, BError> {
        let repo = Repository::open(cwd.clone())?;
        let bin = std::env::current_exe()?;

        Ok(Self { repo, cwd, bin })
    }
}

impl Display for GitRepoContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "Context {{ repo: {:?}, cwd: {:?} }}", self.repo.path(), self.cwd)
    }
}
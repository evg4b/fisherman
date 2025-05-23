use crate::context::Context;
use anyhow::{bail, Result};
use git2::Repository;
use std::fmt::Display;
use std::fs;
use std::path::{Path, PathBuf};

pub struct GitRepoContext {
    repo: Repository,
    cwd: PathBuf,
    bin: PathBuf,
    message_file: Option<PathBuf>,
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

    fn current_branch(&self) -> Result<String> {
        let head = self.repo.head()?;
        Ok(head.shorthand().unwrap_or("HEAD").to_string())
    }

    fn commit_msg(&self) -> Result<String> {
        if let Some(message_file) = self.message_file.as_ref() {
            let message = fs::read_to_string(message_file)?;
            return Ok(message);
        }

        bail!("This action does not support commit message extraction");
    }

    fn set_commit_msg_path(&mut self, message_file: PathBuf) {
        self.message_file = Some(message_file);
    }
}

impl GitRepoContext {
    pub(crate) fn new(cwd: PathBuf) -> Result<Self> {
        let repo = Repository::open(cwd.clone())?;
        let bin = std::env::current_exe()?;

        Ok(Self { repo, cwd, bin, message_file: None })
    }
}

impl Display for GitRepoContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Context {{ repo: {:?}, cwd: {:?} }}",
            self.repo.path(),
            self.cwd
        )
    }
}
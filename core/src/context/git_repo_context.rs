use crate::configuration::Configuration;
use crate::context::variables::extract_variables;
use crate::context::Context;
use anyhow::{bail, Result};
use git2::Repository;
use std::collections::HashMap;
use std::fmt::Display;
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::Mutex;
use figment::providers::Serialized;

pub struct GitRepoContext {
    repo: Mutex<Repository>,
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
        let repo = self.repo.lock().unwrap();
        let head = repo.head()?;
        Ok(head.shorthand().unwrap_or("HEAD").to_string())
    }

    fn commit_msg(&self) -> Result<String> {
        if let Some(message_file) = self.message_file.as_ref() {
            let message = fs::read_to_string(message_file)?;
            // Git adds trailing newlines to commit messages, strip them for validation
            return Ok(message.trim_end().to_string());
        }

        bail!("Commit message not available for this hook");
    }

    fn set_commit_msg_path(&mut self, message_file: PathBuf) {
        self.message_file = Some(message_file);
    }

    fn configuration(&self) -> Result<Configuration> {
        Configuration::load(self.repo_path())
    }

    fn variables(&self, additional: &[String]) -> Result<HashMap<String, String>> {
        let mut variables = additional.to_vec();
        variables.extend(self.configuration()?.extract);
        extract_variables(self, &variables)
    }

    fn staged_files(&self) -> Result<Vec<PathBuf>> {
        let repo = self.repo.lock().unwrap();
        let mut status_options = git2::StatusOptions::new();
        status_options.include_untracked(false);
        status_options.show(git2::StatusShow::Index);
        let statuses = repo.statuses(Some(&mut status_options))?;

        let mut files = Vec::new();
        for entry in statuses.iter() {
            if let Some(path) = entry.path() {
                files.push(PathBuf::from(path));
            }
        }
        Ok(files)
    }
}

impl GitRepoContext {
    pub fn new(cwd: PathBuf) -> Result<Self> {
        let repo = Repository::open(&cwd)?;
        let bin = std::env::current_exe()?;

        Ok(Self {
            repo: Mutex::new(repo),
            cwd,
            bin,
            message_file: None,
        })
    }
}

impl Display for GitRepoContext {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Context {{ repo: {:?}, cwd: {:?} }}",
            self.repo.lock().unwrap().path(),
            self.cwd
        )
    }
}

#[cfg(test)]
pub mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_context_initialization() {
        let temp_dir = tempdir().unwrap();
        // Initialize a real git repo in the temp dir
        Repository::init(temp_dir.path()).unwrap();

        let ctx = GitRepoContext::new(temp_dir.path().to_path_buf()).unwrap();
        assert_eq!(ctx.repo_path(), temp_dir.path());
        assert_eq!(ctx.cwd, temp_dir.path());
        assert_eq!(ctx.bin, std::env::current_exe().unwrap());
    }
}
use crate::configuration::Configuration;
use crate::context::variables::extract_variables;
use crate::context::{Context, DiffLine};
use anyhow::{anyhow, bail, Result};
use git2::Repository;
use std::collections::HashMap;
use std::fmt::Display;
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::Mutex;

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
        let repo = self
            .repo
            .lock()
            .map_err(|e| anyhow!("Failed to lock repository: {}", e))?;
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
        let repo = self
            .repo
            .lock()
            .map_err(|e| anyhow!("Failed to lock repository: {}", e))?;
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

    fn staged_diff(&self, path: &Path) -> Result<Vec<DiffLine>> {
        let repo = self
            .repo
            .lock()
            .map_err(|e| anyhow!("Failed to lock repository: {}", e))?;
        let head = match repo.head() {
            Ok(reference) => Some(reference.peel_to_tree()?),
            Err(_) => None,
        };

        let mut diff_options = git2::DiffOptions::new();
        diff_options.pathspec(path);

        let diff = repo.diff_tree_to_index(head.as_ref(), None, Some(&mut diff_options))?;
        let mut diff_lines = Vec::new();

        diff.foreach(
            &mut |_, _| true,
            None,
            None,
            Some(&mut |_, _, line| {
                let content = std::str::from_utf8(line.content())
                    .map(|s| s.trim_end().to_string())
                    .unwrap_or_default();

                match line.origin() {
                    '+' => diff_lines.push(DiffLine::Added(content)),
                    '-' => diff_lines.push(DiffLine::Deleted(content)),
                    _ => {}
                }
                true
            }),
        )?;

        Ok(diff_lines)
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
        let repo_path = self.repo.lock().map(|r| r.path().to_path_buf()).ok();
        write!(
            f,
            "Context {{ repo: {:?}, cwd: {:?} }}",
            repo_path.map(|p| p.display().to_string()).unwrap_or_else(|| "<locked>".to_string()),
            self.cwd
        )
    }
}

#[cfg(test)]
pub mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn test_context_initialization() -> Result<()> {
        let temp_dir = tempdir()?;
        // Initialize a real git repo in the temp dir
        Repository::init(temp_dir.path())?;

        let ctx = GitRepoContext::new(temp_dir.path().to_path_buf())?;
        assert_eq!(ctx.repo_path(), temp_dir.path());
        assert_eq!(ctx.cwd, temp_dir.path());
        assert_eq!(ctx.bin, std::env::current_exe()?);
        Ok(())
    }
}
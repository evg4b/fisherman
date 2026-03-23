use crate::configuration::Configuration;
use crate::context::variables::extract_variables;
use crate::context::{Context, DiffLine};
use anyhow::{anyhow, bail, Result};
use git2::Repository;
use std::collections::HashMap;
use std::fmt::Display;
use std::fs;
use std::path::{Path, PathBuf};
use std::sync::{Arc, Mutex};

pub struct GitRepoContext {
    configuration: Arc<Configuration>,
    repo: Mutex<Repository>,
    cwd: PathBuf,
    bin: PathBuf,
    message_file: Option<PathBuf>,
    extract: Option<Vec<String>>,
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

    fn configuration(&self) -> Arc<Configuration> {
        self.configuration.clone()
    }

    fn extend(&mut self, extract: &[String]) -> Result<Box<dyn Context>> {
        Ok(Box::new(GitRepoContext {
            configuration: self.configuration.clone(),
            repo: Mutex::new(Repository::open(self.repo_path())?),
            cwd: self.cwd.clone(),
            bin: self.bin.clone(),
            message_file: self.message_file.clone(),
            extract: Option::from(extract.to_vec()),
        }))
    }

    fn variables(&self) -> Result<HashMap<String, String>> {
        match &self.extract {
            None => self.compute_variables(&[]),
            Some(additional) => self.compute_variables(additional),
        }
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
            configuration: Arc::new(Configuration::load(cwd.as_path())?),
            repo: Mutex::new(repo),
            cwd,
            bin,
            message_file: None,
            extract: None,
        })
    }

    fn compute_variables(&self, additional: &[String]) -> Result<HashMap<String, String>> {
        let mut variables = additional.to_vec();
        variables.extend(self.configuration().extract.clone());
        extract_variables(self, &variables)
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

    fn init_ctx(path: &std::path::Path) -> Result<GitRepoContext> {
        Repository::init(path)?;
        GitRepoContext::new(path.to_path_buf())
    }

    #[test]
    fn test_context_initialization() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        assert_eq!(ctx.repo_path(), temp_dir.path());
        assert_eq!(ctx.cwd, temp_dir.path());
        assert_eq!(ctx.bin, std::env::current_exe()?);
        Ok(())
    }

    #[test]
    fn test_hooks_dir() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        assert_eq!(ctx.hooks_dir(), temp_dir.path().join(".git/hooks"));
        Ok(())
    }

    #[test]
    fn test_bin() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        assert_eq!(ctx.bin(), std::env::current_exe()?.as_path());
        Ok(())
    }

    #[test]
    fn test_commit_msg_no_file_returns_error() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        assert!(ctx.commit_msg().is_err());
        Ok(())
    }

    #[test]
    fn test_commit_msg_with_file() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx(temp_dir.path())?;
        let msg_file = temp_dir.path().join("COMMIT_EDITMSG");
        fs::write(&msg_file, "feat: test commit\n")?;
        ctx.set_commit_msg_path(msg_file);
        assert_eq!(ctx.commit_msg()?, "feat: test commit");
        Ok(())
    }

    #[test]
    fn test_commit_msg_strips_trailing_newlines() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx(temp_dir.path())?;
        let msg_file = temp_dir.path().join("COMMIT_EDITMSG");
        fs::write(&msg_file, "fix: a bug\n\n")?;
        ctx.set_commit_msg_path(msg_file);
        assert_eq!(ctx.commit_msg()?, "fix: a bug");
        Ok(())
    }

    #[test]
    fn test_set_commit_msg_path() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx(temp_dir.path())?;
        assert!(ctx.commit_msg().is_err());
        let msg_file = temp_dir.path().join("MSG");
        fs::write(&msg_file, "chore: setup\n")?;
        ctx.set_commit_msg_path(msg_file);
        assert!(ctx.commit_msg().is_ok());
        Ok(())
    }

    #[test]
    fn test_current_branch() -> Result<()> {
        let temp_dir = tempdir()?;
        let repo = Repository::init(temp_dir.path())?;
        let sig = git2::Signature::now("Test", "test@test.com")?;
        let tree_id = {
            let mut index = repo.index()?;
            index.write_tree()?
        };
        let tree = repo.find_tree(tree_id)?;
        repo.commit(Some("HEAD"), &sig, &sig, "initial", &tree, &[])?;

        let ctx = GitRepoContext::new(temp_dir.path().to_path_buf())?;
        let branch = ctx.current_branch()?;
        assert!(!branch.is_empty());
        Ok(())
    }

    #[test]
    fn test_staged_files() -> Result<()> {
        let temp_dir = tempdir()?;
        let repo = Repository::init(temp_dir.path())?;
        let file_path = temp_dir.path().join("test.txt");
        fs::write(&file_path, "content")?;
        let mut index = repo.index()?;
        index.add_path(Path::new("test.txt"))?;
        index.write()?;

        let ctx = GitRepoContext::new(temp_dir.path().to_path_buf())?;
        let files = ctx.staged_files()?;
        assert!(files.contains(&PathBuf::from("test.txt")));
        Ok(())
    }

    #[test]
    fn test_staged_files_empty_repo() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        let files = ctx.staged_files()?;
        assert!(files.is_empty());
        Ok(())
    }

    #[test]
    fn test_staged_diff_new_file() -> Result<()> {
        let temp_dir = tempdir()?;
        let repo = Repository::init(temp_dir.path())?;
        let file_path = temp_dir.path().join("test.txt");
        fs::write(&file_path, "hello\nworld\n")?;
        let mut index = repo.index()?;
        index.add_path(Path::new("test.txt"))?;
        index.write()?;

        let ctx = GitRepoContext::new(temp_dir.path().to_path_buf())?;
        let diff = ctx.staged_diff(Path::new("test.txt"))?;
        let added: Vec<_> = diff.iter().filter(|l| matches!(l, DiffLine::Added(_))).collect();
        assert!(!added.is_empty());
        Ok(())
    }

    #[test]
    fn test_staged_diff_empty_repo() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        let diff = ctx.staged_diff(Path::new("nonexistent.txt"))?;
        assert!(diff.is_empty());
        Ok(())
    }

    #[test]
    fn test_variables_empty_extract() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        let vars = ctx.compute_variables(&[])?;
        assert!(vars.is_empty());
        Ok(())
    }

    fn init_ctx_with_commit(path: &std::path::Path) -> Result<GitRepoContext> {
        let repo = Repository::init(path)?;
        let sig = git2::Signature::now("Test", "test@test.com")?;
        let tree_id = repo.index()?.write_tree()?;
        let tree = repo.find_tree(tree_id)?;
        repo.commit(Some("HEAD"), &sig, &sig, "initial", &tree, &[])?;
        GitRepoContext::new(path.to_path_buf())
    }

    #[test]
    fn test_variables_none_extract_returns_empty() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        // extract is None by default after new()
        let vars = ctx.variables()?;
        assert!(vars.is_empty());
        Ok(())
    }

    #[test]
    fn test_variables_with_extract_from_branch() -> Result<()> {
        let temp_dir = tempdir()?;
        let repo = Repository::init(temp_dir.path())?;
        let sig = git2::Signature::now("Test", "test@test.com")?;
        let tree_id = repo.index()?.write_tree()?;
        let tree = repo.find_tree(tree_id)?;
        repo.commit(Some("HEAD"), &sig, &sig, "initial", &tree, &[])?;
        // rename branch to feat/my-ticket
        repo.branch("feat/my-ticket", &repo.head()?.peel_to_commit()?, false)?;
        repo.set_head("refs/heads/feat/my-ticket")?;

        let mut ctx = GitRepoContext::new(temp_dir.path().to_path_buf())?;
        let extended = ctx.extend(&["branch:^(?P<Type>feat|fix)/(?P<Name>.+)".to_string()])?;
        let vars = extended.variables()?;

        assert_eq!(vars.get("Type"), Some(&"feat".to_string()));
        assert_eq!(vars.get("Name"), Some(&"my-ticket".to_string()));
        Ok(())
    }

    #[test]
    fn test_extend_creates_new_context_with_extract() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx_with_commit(temp_dir.path())?;
        let extract = vec!["branch:^(?P<Name>.+)".to_string()];

        let extended = ctx.extend(&extract)?;
        let vars = extended.variables()?;

        assert!(vars.contains_key("Name"));
        Ok(())
    }

    #[test]
    fn test_extend_inherits_message_file() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx_with_commit(temp_dir.path())?;
        let msg_file = temp_dir.path().join("COMMIT_EDITMSG");
        fs::write(&msg_file, "feat: something\n")?;
        ctx.set_commit_msg_path(msg_file);

        let extended = ctx.extend(&[])?;
        assert_eq!(extended.commit_msg()?, "feat: something");
        Ok(())
    }

    #[test]
    fn test_extend_fails_on_invalid_repo() -> Result<()> {
        let temp_dir = tempdir()?;
        let mut ctx = init_ctx(temp_dir.path())?;
        // Remove .git to make Repository::open fail
        std::fs::remove_dir_all(temp_dir.path().join(".git"))?;

        let result = ctx.extend(&[]);
        assert!(result.is_err());
        Ok(())
    }

    #[test]
    fn test_display_format() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        let display = format!("{}", ctx);
        assert!(display.contains("Context"));
        Ok(())
    }

    #[test]
    fn test_configuration_empty_dir() -> Result<()> {
        let temp_dir = tempdir()?;
        let ctx = init_ctx(temp_dir.path())?;
        let config = ctx.configuration();
        assert!(config.hooks.is_empty());
        Ok(())
    }
}
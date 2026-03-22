mod context_trait;
mod git_repo_context;
mod variables;

pub use context_trait::{Context, DiffLine, MockContext};
pub use git_repo_context::GitRepoContext;

mod context;
mod git_repo_context;

pub use context::{Context};
pub use git_repo_context::GitRepoContext;

#[cfg(test)]
pub use context::{MockContext};
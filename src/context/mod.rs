mod context_trait;
mod git_repo_context;
mod variables;

pub use context_trait::Context;
pub use git_repo_context::GitRepoContext;

#[cfg(test)]
pub use context_trait::MockContext;

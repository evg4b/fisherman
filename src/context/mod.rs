#[allow(clippy::module_inception)]
mod context;
mod git_repo_context;
mod variables;

pub use context::Context;
pub use git_repo_context::GitRepoContext;

#[cfg(test)]
pub use context::MockContext;

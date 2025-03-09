mod errors;
mod files;
mod git_hook;

pub use git_hook::GitHook;
pub use files::{build_hook_content, override_hook, write_hook};
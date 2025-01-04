pub(crate) mod errors;
pub(crate) mod files;

use clap::ValueEnum;
use serde::Deserialize;

const APPLYPATCH_MSG: &str = "applypatch-msg";
const COMMIT_MSG: &str = "commit-msg";
const FSMONITOR_WATCHMAN: &str = "fsmonitor-watchman";
const POST_UPDATE: &str = "post-update";
const PRE_APPLY_PATCH: &str = "pre-applypatch";
const PRE_COMMIT: &str = "pre-commit";
const PRE_MERGE_COMMIT: &str = "pre-merge-commit";
const PRE_PUSH: &str = "pre-push";
const PRE_REBASE: &str = "pre-rebase";
const PRE_RECEIVE: &str = "pre-receive";
const PREPARE_COMMIT_MSG: &str = "prepare-commit-msg";
const PUSH_TO_CHECKOUT: &str = "push-to-checkout";
const SENDEMAIL_VALIDATE: &str = "sendemail-validate";
const UPDATE: &str = "update";

#[derive(Debug, Deserialize, Hash, Eq, PartialEq, Copy, Clone, ValueEnum)]
pub(crate) enum GitHook {
    #[serde(rename = "applypatch-msg")]
    ApplypatchMsg,
    #[serde(rename = "commit-msg")]
    CommitMsg,
    #[serde(rename = "fsmonitor-watchman")]
    FsmonitorWatchman,
    #[serde(rename = "post-update")]
    PostUpdate,
    #[serde(rename = "pre-applypatch")]
    PreApplypatch,
    #[serde(rename = "pre-commit")]
    PreCommit,
    #[serde(rename = "pre-merge-commit")]
    PreMergeCommit,
    #[serde(rename = "pre-push")]
    PrePush,
    #[serde(rename = "pre-rebase")]
    PreRebase,
    #[serde(rename = "pre-receive")]
    PreReceive,
    #[serde(rename = "prepare-commit-msg")]
    PrepareCommitMsg,
    #[serde(rename = "push-to-checkout")]
    PushToCheckout,
    #[serde(rename = "sendemail-validate")]
    SendemailValidate,
    #[serde(rename = "update")]
    Update,
}

impl GitHook {
    // Return all the git hooks
    pub fn all() -> Vec<GitHook> {
        vec![
            GitHook::ApplypatchMsg,
            GitHook::CommitMsg,
            GitHook::FsmonitorWatchman,
            GitHook::PostUpdate,
            GitHook::PreApplypatch,
            GitHook::PreCommit,
            GitHook::PreMergeCommit,
            GitHook::PrePush,
            GitHook::PreRebase,
            GitHook::PreReceive,
            GitHook::PrepareCommitMsg,
            GitHook::PushToCheckout,
            GitHook::SendemailValidate,
            GitHook::Update,
        ]
    }

    // Convert the git hook enum to a string slice
    pub fn as_str(&self) -> &'static str {
        match self {
            GitHook::ApplypatchMsg => APPLYPATCH_MSG,
            GitHook::CommitMsg => COMMIT_MSG,
            GitHook::FsmonitorWatchman => FSMONITOR_WATCHMAN,
            GitHook::PostUpdate => POST_UPDATE,
            GitHook::PreApplypatch => PRE_APPLY_PATCH,
            GitHook::PreCommit => PRE_COMMIT,
            GitHook::PreMergeCommit => PRE_MERGE_COMMIT,
            GitHook::PrePush => PRE_PUSH,
            GitHook::PreRebase => PRE_REBASE,
            GitHook::PreReceive => PRE_RECEIVE,
            GitHook::PrepareCommitMsg => PREPARE_COMMIT_MSG,
            GitHook::PushToCheckout => PUSH_TO_CHECKOUT,
            GitHook::SendemailValidate => SENDEMAIL_VALIDATE,
            GitHook::Update => UPDATE,
        }
    }
}

impl std::fmt::Display for GitHook {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}", self.as_str())
    }
}

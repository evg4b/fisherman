use clap::ValueEnum;
use std::os::unix::fs::PermissionsExt;
use std::path::PathBuf;
use std::{fs, io};

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

#[derive(Copy, Clone, PartialEq, Eq, PartialOrd, ValueEnum, Ord, Debug)]
pub(crate) enum GitHook {
    ApplypatchMsg,
    CommitMsg,
    FsmonitorWatchman,
    PostUpdate,
    PreApplypatch,
    PreCommit,
    PreMergeCommit,
    PrePush,
    PreRebase,
    PreReceive,
    PrepareCommitMsg,
    PushToCheckout,
    SendemailValidate,
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

pub(crate) fn write_hook(path: &PathBuf, hook: GitHook, content: String) -> io::Result<()> {
    let hook_path = &path.join(".git/hooks").join(hook.as_str());
    fs::write(hook_path, content)?;
    fs::set_permissions(hook_path, fs::Permissions::from_mode(0o700))
}

pub(crate) fn build_hook_content(bin: &PathBuf, hook_name: GitHook) -> String {
    format!("#!/bin/sh\n{} handle {}\n", bin.display(), hook_name)
}

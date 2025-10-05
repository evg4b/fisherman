use crate::hooks::errors::HookError;
use anyhow::{bail, Result};
use clap::ValueEnum;
use serde::Deserialize;
use std::fs;
use std::os::unix::fs::PermissionsExt;

const APPLYPATCH_MSG: &str = "applypatch-msg";
const PRE_APPLYPATCH: &str = "pre-applypatch";
const POST_APPLYPATCH: &str = "post-applypatch";
const PRE_COMMIT: &str = "pre-commit";
const PRE_MERGE_COMMIT: &str = "pre-merge-commit";
const PREPARE_COMMIT_MSG: &str = "prepare-commit-msg";
const COMMIT_MSG: &str = "commit-msg";
const POST_COMMIT: &str = "post-commit";
const PRE_REBASE: &str = "pre-rebase";
const POST_CHECKOUT: &str = "post-checkout";
const POST_MERGE: &str = "post-merge";
const PRE_PUSH: &str = "pre-push";
const PRE_RECEIVE: &str = "pre-receive";
const UPDATE: &str = "update";
const PROC_RECEIVE: &str = "proc-receive";
const POST_RECEIVE: &str = "post-receive";
const POST_UPDATE: &str = "post-update";
const REFERENCE_TRANSACTION: &str = "reference-transaction";
const PUSH_TO_CHECKOUT: &str = "push-to-checkout";
const PRE_AUTO_GC: &str = "pre-auto-gc";
const POST_REWRITE: &str = "post-rewrite";
const SENDEMAIL_VALIDATE: &str = "sendemail-validate";
const FSMONITOR_WATCHMAN: &str = "fsmonitor-watchman";
const P4_CHANGELIST: &str = "p4-changelist";
const P4_PREPARE_CHANGELIST: &str = "p4-prepare-changelist";
const P4_POST_CHANGELIST: &str = "p4-post-changelist";
const P4_PRE_SUBMIT: &str = "p4-pre-submit";
const POST_INDEX_CHANGE: &str = "post-index-change";

#[derive(Debug, Deserialize, Hash, Eq, PartialEq, Copy, Clone, ValueEnum)]
pub enum GitHook {
    #[serde(rename = "applypatch-msg")]
    ApplypatchMsg,
    #[serde(rename = "pre-applypatch")]
    PreApplypatch,
    #[serde(rename = "post-applypatch")]
    PostApplypatch,
    #[serde(rename = "pre-commit")]
    PreCommit,
    #[serde(rename = "pre-merge-commit")]
    PreMergeCommit,
    #[serde(rename = "prepare-commit-msg")]
    PrepareCommitMsg,
    #[serde(rename = "commit-msg")]
    CommitMsg,
    #[serde(rename = "post-commit")]
    PostCommit,
    #[serde(rename = "pre-rebase")]
    PreRebase,
    #[serde(rename = "post-checkout")]
    PostCheckout,
    #[serde(rename = "post-merge")]
    PostMerge,
    #[serde(rename = "pre-push")]
    PrePush,
    #[serde(rename = "pre-receive")]
    PreReceive,
    #[serde(rename = "update")]
    Update,
    #[serde(rename = "proc-receive")]
    ProcReceive,
    #[serde(rename = "post-receive")]
    PostReceive,
    #[serde(rename = "post-update")]
    PostUpdate,
    #[serde(rename = "reference-transaction")]
    ReferenceTransaction,
    #[serde(rename = "push-to-checkout")]
    PushToCheckout,
    #[serde(rename = "pre-auto-gc")]
    PreAutoGc,
    #[serde(rename = "post-rewrite")]
    PostRewrite,
    #[serde(rename = "sendemail-validate")]
    SendemailValidate,
    #[serde(rename = "fsmonitor-watchman")]
    FsmonitorWatchman,
    #[serde(rename = "p4-changelist")]
    P4Changelist,
    #[serde(rename = "p4-prepare-changelist")]
    P4PrepareChangelist,
    #[serde(rename = "p4-post-changelist")]
    P4PostChangelist,
    #[serde(rename = "p4-pre-submit")]
    P4PreSubmit,
    #[serde(rename = "post-index-change")]
    PostIndexChange,
}

impl GitHook {
    // Convert the git hook enum to a string slice
    pub fn as_str(&self) -> &'static str {
        match self {
            GitHook::ApplypatchMsg => APPLYPATCH_MSG,
            GitHook::PreApplypatch => PRE_APPLYPATCH,
            GitHook::PostApplypatch => POST_APPLYPATCH,
            GitHook::PreCommit => PRE_COMMIT,
            GitHook::PreMergeCommit => PRE_MERGE_COMMIT,
            GitHook::PrepareCommitMsg => PREPARE_COMMIT_MSG,
            GitHook::CommitMsg => COMMIT_MSG,
            GitHook::PostCommit => POST_COMMIT,
            GitHook::PreRebase => PRE_REBASE,
            GitHook::PostCheckout => POST_CHECKOUT,
            GitHook::PostMerge => POST_MERGE,
            GitHook::PrePush => PRE_PUSH,
            GitHook::PreReceive => PRE_RECEIVE,
            GitHook::Update => UPDATE,
            GitHook::ProcReceive => PROC_RECEIVE,
            GitHook::PostReceive => POST_RECEIVE,
            GitHook::PostUpdate => POST_UPDATE,
            GitHook::ReferenceTransaction => REFERENCE_TRANSACTION,
            GitHook::PushToCheckout => PUSH_TO_CHECKOUT,
            GitHook::PreAutoGc => PRE_AUTO_GC,
            GitHook::PostRewrite => POST_REWRITE,
            GitHook::SendemailValidate => SENDEMAIL_VALIDATE,
            GitHook::FsmonitorWatchman => FSMONITOR_WATCHMAN,
            GitHook::P4Changelist => P4_CHANGELIST,
            GitHook::P4PrepareChangelist => P4_PREPARE_CHANGELIST,
            GitHook::P4PostChangelist => P4_POST_CHANGELIST,
            GitHook::P4PreSubmit => P4_PRE_SUBMIT,
            GitHook::PostIndexChange => POST_INDEX_CHANGE,
        }
    }

    pub fn install(&self, context: &impl crate::context::Context, force: bool) -> Result<()> {
        let hook_path = &context.hooks_dir().join(self.as_str());
        let hook_exists = hook_path.exists();
        if hook_exists && !force {
            bail!(HookError::AlreadyExists {
                name: self.as_str(),
                hook: hook_path.to_owned()
            });
        }

        if hook_exists {
            fs::copy(hook_path, hook_path.with_extension("bkp"))?;
            fs::remove_file(hook_path)?;
        }

        fs::write(hook_path, self.content(context))?;
        fs::set_permissions(hook_path, fs::Permissions::from_mode(0o755))?;

        Ok(())
    }

    fn content(&self, context: &impl crate::context::Context) -> String {
        let bin = context.bin().display();
        match self {
            GitHook::CommitMsg => format!("#!/bin/sh\n{} handle {} $@\n", bin, self),
            _ => format!("#!/bin/sh\n{} handle {}\n", bin, self),
        }
    }
}

impl std::fmt::Display for GitHook {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}", self.as_str())
    }
}

#[cfg(test)]
mod test_hook_install {
    use super::*;
    use crate::context::MockContext;
    use crate::hooks::GitHook;
    use assert2::assert;
    use rstest::*;
    use std::path::PathBuf;
    use tempdir::TempDir;

    #[rstest]
    #[case(APPLYPATCH_MSG)]
    #[case(PRE_APPLYPATCH)]
    #[case(POST_APPLYPATCH)]
    #[case(PRE_COMMIT)]
    #[case(PRE_MERGE_COMMIT)]
    #[case(PREPARE_COMMIT_MSG)]
    #[case(COMMIT_MSG)]
    #[case(POST_COMMIT)]
    #[case(PRE_REBASE)]
    #[case(POST_CHECKOUT)]
    #[case(POST_MERGE)]
    #[case(PRE_PUSH)]
    #[case(PRE_RECEIVE)]
    #[case(UPDATE)]
    #[case(PROC_RECEIVE)]
    #[case(POST_RECEIVE)]
    #[case(POST_UPDATE)]
    #[case(REFERENCE_TRANSACTION)]
    #[case(PUSH_TO_CHECKOUT)]
    #[case(PRE_AUTO_GC)]
    #[case(POST_REWRITE)]
    #[case(SENDEMAIL_VALIDATE)]
    #[case(FSMONITOR_WATCHMAN)]
    #[case(P4_CHANGELIST)]
    #[case(P4_PREPARE_CHANGELIST)]
    #[case(P4_POST_CHANGELIST)]
    #[case(P4_PRE_SUBMIT)]
    #[case(POST_INDEX_CHANGE)]
    fn install_test(#[case] hook_name: &str) {
        let hook = GitHook::from_str(hook_name, false).unwrap();
        let dir = TempDir::new(format!("test_install_{}", hook_name).as_str()).unwrap();

        let mut ctx = MockContext::new();
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        hook.install(&ctx, false).unwrap();

        let hook_path = dir.into_path().join(hook.to_string());

        assert!(hook_path.exists());
        assert!(hook_path.is_file());
        assert!(fs::read_to_string(hook_path).unwrap() == hook.content(&ctx));
    }

    #[rstest]
    #[case(APPLYPATCH_MSG)]
    #[case(PRE_APPLYPATCH)]
    #[case(POST_APPLYPATCH)]
    #[case(PRE_COMMIT)]
    #[case(PRE_MERGE_COMMIT)]
    #[case(PREPARE_COMMIT_MSG)]
    #[case(COMMIT_MSG)]
    #[case(POST_COMMIT)]
    #[case(PRE_REBASE)]
    #[case(POST_CHECKOUT)]
    #[case(POST_MERGE)]
    #[case(PRE_PUSH)]
    #[case(PRE_RECEIVE)]
    #[case(UPDATE)]
    #[case(PROC_RECEIVE)]
    #[case(POST_RECEIVE)]
    #[case(POST_UPDATE)]
    #[case(REFERENCE_TRANSACTION)]
    #[case(PUSH_TO_CHECKOUT)]
    #[case(PRE_AUTO_GC)]
    #[case(POST_REWRITE)]
    #[case(SENDEMAIL_VALIDATE)]
    #[case(FSMONITOR_WATCHMAN)]
    #[case(P4_CHANGELIST)]
    #[case(P4_PREPARE_CHANGELIST)]
    #[case(P4_POST_CHANGELIST)]
    #[case(P4_PRE_SUBMIT)]
    #[case(POST_INDEX_CHANGE)]
    fn install_force_test(#[case] hook_name: &str) {
        let hook = GitHook::from_str(hook_name, false).unwrap();
        let dir = TempDir::new(format!("test_install_force_{}", hook_name).as_str()).unwrap();

        let mut ctx = MockContext::new();
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        let original_hook_content = format!("test {}", hook_name);
        fs::write(dir.path().join(hook.to_string()), &original_hook_content).unwrap();

        hook.install(&ctx, true).unwrap();

        let hook_path = dir.into_path().join(hook.to_string());
        let hook_bkp_path = hook_path.with_extension("bkp");

        assert!(hook_path.exists());
        assert!(hook_path.is_file());
        assert!(fs::read_to_string(hook_path).unwrap() == hook.content(&ctx));

        assert!(hook_bkp_path.exists());
        assert!(hook_bkp_path.is_file());
        assert!(fs::read_to_string(hook_bkp_path).unwrap() == original_hook_content);
    }
}

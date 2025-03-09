use crate::common::BError;
use crate::err;
use crate::hooks::errors::HookError;
use clap::ValueEnum;
use serde::Deserialize;
use std::fs;
use std::os::unix::fs::PermissionsExt;

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
pub enum GitHook {
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

    pub fn install(
        &self,
        context: &impl crate::context::Context,
        force: bool,
    ) -> Result<(), BError> {
        let hook_path = &context.hooks_dir().join(self.as_str());
        let hook_exists = hook_path.exists();
        if hook_exists && !force {
            err!(HookError::AlreadyExists {
                name: self.as_str(),
                hook: hook_path.clone()
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
        format!("#!/bin/sh\n{} handle {}\n", context.bin().display(), self)
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
    use assertor::{EqualityAssertion, assert_that};
    use rstest::*;
    use std::path::PathBuf;
    use tempdir::TempDir;

    #[rstest]
    #[case(APPLYPATCH_MSG)]
    #[case(COMMIT_MSG)]
    #[case(FSMONITOR_WATCHMAN)]
    #[case(POST_UPDATE)]
    #[case(PRE_APPLY_PATCH)]
    #[case(PRE_COMMIT)]
    #[case(PRE_MERGE_COMMIT)]
    #[case(PRE_PUSH)]
    #[case(PRE_REBASE)]
    #[case(PRE_RECEIVE)]
    #[case(PREPARE_COMMIT_MSG)]
    #[case(PUSH_TO_CHECKOUT)]
    #[case(SENDEMAIL_VALIDATE)]
    #[case(UPDATE)]
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
        assert_that!(fs::read_to_string(hook_path).unwrap())
            .is_equal_to(hook.content(&ctx));
    }

    #[rstest]
    #[case(APPLYPATCH_MSG)]
    #[case(COMMIT_MSG)]
    #[case(FSMONITOR_WATCHMAN)]
    #[case(POST_UPDATE)]
    #[case(PRE_APPLY_PATCH)]
    #[case(PRE_COMMIT)]
    #[case(PRE_MERGE_COMMIT)]
    #[case(PRE_PUSH)]
    #[case(PRE_REBASE)]
    #[case(PRE_RECEIVE)]
    #[case(PREPARE_COMMIT_MSG)]
    #[case(PUSH_TO_CHECKOUT)]
    #[case(SENDEMAIL_VALIDATE)]
    #[case(UPDATE)]
    fn install_force_test(#[case] hook_name: &str) {
        let hook = GitHook::from_str(hook_name, false).unwrap();
        let dir = TempDir::new(format!("test_install_force_{}", hook_name).as_str()).unwrap();

        let mut ctx = MockContext::new();
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        let original_hook_content = format!("test {}", hook_name);
        fs::write(
            dir.path().join(hook.to_string()),
            original_hook_content.clone(),
        )
        .unwrap();

        hook.install(&ctx, true).unwrap();

        let hook_path = dir.into_path().join(hook.to_string());
        let hook_bkp_path = hook_path.with_extension("bkp");

        assert!(hook_path.exists());
        assert!(hook_path.is_file());
        assert_that!(fs::read_to_string(hook_path).unwrap())
            .is_equal_to(hook.content(&ctx));

        assert!(hook_bkp_path.exists());
        assert!(hook_bkp_path.is_file());
        assert_that!(fs::read_to_string(hook_bkp_path).unwrap())
            .is_equal_to(original_hook_content);
    }
}

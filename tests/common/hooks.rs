#[derive(Debug, Copy, Clone, PartialEq, Eq, Hash)]
pub enum Hook {
    ApplypatchMsg,
    PreApplypatch,
    PostApplypatch,
    PreCommit,
    PreMergeCommit,
    PrepareCommitMsg,
    CommitMsg,
    PostCommit,
    PreRebase,
    PostCheckout,
    PostMerge,
    PrePush,
    PreReceive,
    Update,
    ProcReceive,
    PostReceive,
    PostUpdate,
    ReferenceTransaction,
    PushToCheckout,
    PreAutoGc,
    PostRewrite,
    SendemailValidate,
    FsmonitorWatchman,
    P4Changelist,
    P4PrepareChangelist,
    P4PostChangelist,
    P4PreSubmit,
    PostIndexChange,
}

impl Hook {
    /// Возвращает срез со всеми вариантами enum
    pub fn iter() -> &'static [Hook] {
        &[
            Hook::ApplypatchMsg,
            Hook::PreApplypatch,
            Hook::PostApplypatch,
            Hook::PreCommit,
            Hook::PreMergeCommit,
            Hook::PrepareCommitMsg,
            Hook::CommitMsg,
            Hook::PostCommit,
            Hook::PreRebase,
            Hook::PostCheckout,
            Hook::PostMerge,
            Hook::PrePush,
            Hook::PreReceive,
            Hook::Update,
            Hook::ProcReceive,
            Hook::PostReceive,
            Hook::PostUpdate,
            Hook::ReferenceTransaction,
            Hook::PushToCheckout,
            Hook::PreAutoGc,
            Hook::PostRewrite,
            Hook::SendemailValidate,
            Hook::FsmonitorWatchman,
            Hook::P4Changelist,
            Hook::P4PrepareChangelist,
            Hook::P4PostChangelist,
            Hook::P4PreSubmit,
            Hook::PostIndexChange,
        ]
    }

    /// Возвращает строковое значение, аналог serde(rename)
    pub fn as_str(&self) -> &'static str {
        match self {
            Hook::ApplypatchMsg => "applypatch-msg",
            Hook::PreApplypatch => "pre-applypatch",
            Hook::PostApplypatch => "post-applypatch",
            Hook::PreCommit => "pre-commit",
            Hook::PreMergeCommit => "pre-merge-commit",
            Hook::PrepareCommitMsg => "prepare-commit-msg",
            Hook::CommitMsg => "commit-msg",
            Hook::PostCommit => "post-commit",
            Hook::PreRebase => "pre-rebase",
            Hook::PostCheckout => "post-checkout",
            Hook::PostMerge => "post-merge",
            Hook::PrePush => "pre-push",
            Hook::PreReceive => "pre-receive",
            Hook::Update => "update",
            Hook::ProcReceive => "proc-receive",
            Hook::PostReceive => "post-receive",
            Hook::PostUpdate => "post-update",
            Hook::ReferenceTransaction => "reference-transaction",
            Hook::PushToCheckout => "push-to-checkout",
            Hook::PreAutoGc => "pre-auto-gc",
            Hook::PostRewrite => "post-rewrite",
            Hook::SendemailValidate => "sendemail-validate",
            Hook::FsmonitorWatchman => "fsmonitor-watchman",
            Hook::P4Changelist => "p4-changelist",
            Hook::P4PrepareChangelist => "p4-prepare-changelist",
            Hook::P4PostChangelist => "p4-post-changelist",
            Hook::P4PreSubmit => "p4-pre-submit",
            Hook::PostIndexChange => "post-index-change",
        }
    }
}

impl std::fmt::Display for Hook {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(f, "{}", self.as_str())
    }
}
package constants

var (
	ApplypatchMsgHook        = "applypatch-msg"
	PreApplypatchHook        = "pre-applypatch"
	PostApplypatchHook       = "post-applypatch"
	PreCommitHook            = "pre-commit"
	PreMergeCommitHook       = "pre-merge-commit"
	PrepareCommitMsgHook     = "prepare-commit-msg"
	CommitMsgHook            = "commit-msg"
	PostCommitHook           = "post-commit"
	PreRebaseHook            = "pre-rebase"
	PostCheckoutHook         = "post-checkout"
	PostMergeHook            = "post-merge"
	PrePushHook              = "pre-push"
	PreReceiveHook           = "pre-receive"
	UpdateHook               = "update"
	ProcReceiveHook          = "proc-receive"
	PostReceiveHook          = "post-receive"
	PostUpdateHook           = "post-update"
	ReferenceTransactionHook = "reference-transaction"
	PushToCheckoutHook       = "push-to-checkout"
	PreAutoGcHook            = "pre-auto-gc"
	PostRewriteHook          = "post-rewrite"
	SendemailValidateHook    = "sendemail-validate"
	FsmonitorWatchmanHook    = "fsmonitor-watchman"
	P4ChangelistHook         = "p4-changelist"
	P4PrepareChangelistHook  = "p4-prepare-changelist"
	P4PostChangelistHook     = "p4-post-changelist"
	P4PreSubmitHook          = "p4-pre-submit"
	PostIndexChangeHook      = "post-index-change"
)

// HooksNames is list on supported hooks.
var HooksNames = []string{
	ApplypatchMsgHook,
	PreApplypatchHook,
	PostApplypatchHook,
	PreCommitHook,
	PreMergeCommitHook,
	PrepareCommitMsgHook,
	CommitMsgHook,
	PostCommitHook,
	PreRebaseHook,
	PostCheckoutHook,
	PostMergeHook,
	PrePushHook,
	PreReceiveHook,
	UpdateHook,
	ProcReceiveHook,
	PostReceiveHook,
	PostUpdateHook,
	ReferenceTransactionHook,
	PushToCheckoutHook,
	PreAutoGcHook,
	PostRewriteHook,
	SendemailValidateHook,
	FsmonitorWatchmanHook,
	P4ChangelistHook,
	P4PrepareChangelistHook,
	P4PostChangelistHook,
	P4PreSubmitHook,
	PostIndexChangeHook,
}

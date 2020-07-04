package constants

var (
	// ApplyPatchMsgHook is constant for applypatch-msg hook
	ApplyPatchMsgHook = "applypatch-msg"
	// CommitMsgHook is constant for commit-msg hook
	CommitMsgHook = "commit-msg"
	// FsMonitorWatchmanHook is constant for fsmonitor-watchman hook
	FsMonitorWatchmanHook = "fsmonitor-watchman"
	// PostUpdateHook is constant for post-update hook
	PostUpdateHook = "post-update"
	// PreApplyPatchHook is constant for pre-applypatch hook
	PreApplyPatchHook = "pre-applypatch"
	// PreCommitHook is constant for pre-commit hook
	PreCommitHook = "pre-commit"
	// PrePushHook is constant for pre-push hook
	PrePushHook = "pre-push"
	// PreRebaseHook is constant for pre-rebase hook
	PreRebaseHook = "pre-rebase"
	// PreReceiveHook is constant for pre-receive hook
	PreReceiveHook = "pre-receive"
	// PrepareCommitMsgHook is constant for prepare-commit-msg hook
	PrepareCommitMsgHook = "prepare-commit-msg"
	// UpdateHook is constant for update hook
	UpdateHook = "update"
)

// HooksNames is hook name list
var HooksNames = []string{
	ApplyPatchMsgHook,
	CommitMsgHook,
	FsMonitorWatchmanHook,
	PostUpdateHook,
	PreApplyPatchHook,
	PreCommitHook,
	PrePushHook,
	PreRebaseHook,
	PreReceiveHook,
	PrepareCommitMsgHook,
	UpdateHook,
}

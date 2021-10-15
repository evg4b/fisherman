package constants

var (
	// ApplyPatchMsgHook is applypatch-msg hook constant.
	ApplyPatchMsgHook = "applypatch-msg"
	// CommitMsgHook is commit-msg hook constant.
	CommitMsgHook = "commit-msg"
	// FsMonitorWatchmanHook is fsmonitor-watchman hook constant.
	FsMonitorWatchmanHook = "fsmonitor-watchman"
	// PostUpdateHook is post-update hook constant.
	PostUpdateHook = "post-update"
	// PreApplyPatchHook is pre-applypatch hook constant.
	PreApplyPatchHook = "pre-applypatch"
	// PreCommitHook is pre-commit hook constant.
	PreCommitHook = "pre-commit"
	// PrePushHook is pre-push hook constant.
	PrePushHook = "pre-push"
	// PreRebaseHook is pre-rebase hook constant.
	PreRebaseHook = "pre-rebase"
	// PreReceiveHook is pre-receive hook constant.
	PreReceiveHook = "pre-receive"
	// PrepareCommitMsgHook is prepare-commit-msg hook constant.
	PrepareCommitMsgHook = "prepare-commit-msg"
	// UpdateHook is update hook constant.
	UpdateHook = "update"
)

// HooksNames is list on supported hooks.
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

package constants

var (
	// ApplyPatchMsgHook is applypatch-msg hook constant.
	ApplyPatchMsgHook = "applypatch-msg"
	// CommitMsgHook is commit-msg hook constant.
	CommitMsgHook = "commit-msg"
	// FsMonitorWatchmanHook is fsmonitor-watchman hook constant.
	FsMonitorWatchmanHook = "fsmonitor-watchman"
	// PreApplyPatchHook is pre-applypatch hook constant.
	PreApplyPatchHook = "pre-applypatch"
	// PreCommitHook is pre-commit hook constant.
	PreCommitHook = "pre-commit"
	// PrePushHook is pre-push hook constant.
	PrePushHook = "pre-push"
	// PreRebaseHook is pre-rebase hook constant.
	PreRebaseHook = "pre-rebase"
	// PrepareCommitMsgHook is prepare-commit-msg hook constant.
	PrepareCommitMsgHook = "prepare-commit-msg"
)

// HooksNames is list on supported hooks.
var HooksNames = []string{
	ApplyPatchMsgHook,
	CommitMsgHook,
	FsMonitorWatchmanHook,
	PreApplyPatchHook,
	PreCommitHook,
	PrePushHook,
	PreRebaseHook,
	PrepareCommitMsgHook,
}

package constants

var (
	ApplyPatchMsgHook     = "applypatch-msg"
	CommitMsgHook         = "commit-msg"
	FsMonitorWatchmanHook = "fsmonitor-watchman"
	PostUpdateHook        = "post-update"
	PreApplyPatchHook     = "pre-applypatch"
	PreCommitHook         = "pre-commit"
	PrePushHook           = "pre-push"
	PreRebaseHook         = "pre-rebase"
	PreReceiveHook        = "pre-receive"
	PrepareCommitMsgHook  = "prepare-commit-msg"
	UpdateHook            = "update"
)

var HooksNames = []string{
	// ApplyPatchMsgHook,
	CommitMsgHook,
	// FsMonitorWatchmanHook,
	// PostUpdateHook,
	// PreApplyPatchHook,
	PreCommitHook,
	PrePushHook,
	// PreRebaseHook,
	// PreReceiveHook,
	PrepareCommitMsgHook,
	// UpdateHook,
}

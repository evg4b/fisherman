package constants

const AppName = "fisherman"

var Version = "x.x.x"
var AppConfigNames = []string{".fisherman.yaml", ".fisherman.yml"}

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

const (
	EmailVariable            = "Email"
	UserNameVariable         = "UserName"
	FishermanVersionVariable = "FishermanVersion"
	CwdVariable              = "CWD"
)

const (
	GlobalConfigPath = "GlobalConfigPath"
	LocalConfigPath  = "LocalConfigPath"
	RepoConfigPath   = "RepoConfigPath"
	HookName         = "HookName"
)

package configuration

type HooksConfig struct {
	ApplyPatchMsgHook     ApplyPatchMsgHookConfig     `yaml:"applypatch-msg,omitempty"`
	FsMonitorWatchmanHook FsMonitorWatchmanHookConfig `yaml:"fsmonitor-watchman,omitempty"`
	PostUpdateHook        PostUpdateHookConfig        `yaml:"post-update,omitempty"`
	PreApplyPatchHook     PreApplyPatchHookConfig     `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         PreCommitHookConfig         `yaml:"pre-commit,omitempty"`
	PrePushHook           PrePushHookConfig           `yaml:"pre-push,omitempty"`
	PreRebaseHook         PreRebaseHookConfig         `yaml:"pre-rebase,omitempty"`
	PreReceiveHook        PreReceiveHookConfig        `yaml:"pre-receive,omitempty"`
	UpdateHook            UpdateHookConfig            `yaml:"update,omitempty"`
	CommitMsgHook         CommitMsgHookConfig         `yaml:"commit-msg,omitempty"`
	PrepareCommitMsgHook  PrepareCommitMsgHookConfig  `yaml:"prepare-commit-msg,omitempty"`
}

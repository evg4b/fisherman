package configuration

type HooksConfig struct {
	ApplyPatchMsgHook     *HookConfig `yaml:"applypatch-msg,omitempty"`
	FsMonitorWatchmanHook *HookConfig `yaml:"fsmonitor-watchman,omitempty"`
	PostUpdateHook        *HookConfig `yaml:"post-update,omitempty"`
	PreApplyPatchHook     *HookConfig `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         *HookConfig `yaml:"pre-commit,omitempty"`
	PrePushHook           *HookConfig `yaml:"pre-push,omitempty"`
	PreRebaseHook         *HookConfig `yaml:"pre-rebase,omitempty"`
	PreReceiveHook        *HookConfig `yaml:"pre-receive,omitempty"`
	UpdateHook            *HookConfig `yaml:"update,omitempty"`
	CommitMsgHook         *HookConfig `yaml:"commit-msg,omitempty"`
	PrepareCommitMsgHook  *HookConfig `yaml:"prepare-commit-msg,omitempty"`
}

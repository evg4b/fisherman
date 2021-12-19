package configuration

type HooksConfig struct {
	ApplyPatchMsgHook     *HookConfig `yaml:"applypatch-msg,omitempty"`
	FsMonitorWatchmanHook *HookConfig `yaml:"fsmonitor-watchman,omitempty"`
	PreApplyPatchHook     *HookConfig `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         *HookConfig `yaml:"pre-commit,omitempty"`
	PrePushHook           *HookConfig `yaml:"pre-push,omitempty"`
	PreRebaseHook         *HookConfig `yaml:"pre-rebase,omitempty"`
	CommitMsgHook         *HookConfig `yaml:"commit-msg,omitempty"`
	PrepareCommitMsgHook  *HookConfig `yaml:"prepare-commit-msg,omitempty"`
}

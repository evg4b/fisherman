package config

type HooksConfig struct {
	ApplyPatchMsgHook     *struct{} `yaml:"applypatch-msg,omitempty"`
	CommitMsgHook         *struct{} `yaml:"commit-msg,omitempty"`
	FsMonitorWatchmanHook *struct{} `yaml:"fsmonitor-watchman,omitempty"`
	PostUpdateHook        *struct{} `yaml:"post-update,omitempty"`
	PreApplyPatchHook     *struct{} `yaml:"pre-applypatch,omitempty"`
	PreCommitHook         *struct{} `yaml:"pre-commit,omitempty"`
	PrePushHook           *struct{} `yaml:"pre-push,omitempty"`
	PreRebaseHook         *struct{} `yaml:"pre-rebase,omitempty"`
	PreReceiveHook        *struct{} `yaml:"pre-receive,omitempty"`
	PrepareCommitMsgHook  *struct{} `yaml:"prepare-commit-msg,omitempty"`
	UpdateHook            *struct{} `yaml:"update,omitempty"`
}

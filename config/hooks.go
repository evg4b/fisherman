package config

type HooksConfig struct {
	ApplyPatchMsgHook     string `yaml:"applypatch-msg",omitempty`
	CommitMsgHook         string `yaml:"commit-msg",omitempty`
	FsMonitorWatchmanHook string `yaml:"fsmonitor-watchman",omitempty`
	PostUpdateHook        string `yaml:"post-update",omitempty`
	PreApplyPatchHook     string `yaml:"pre-applypatch",omitempty`
	PreCommitHook         string `yaml:"pre-commit",omitempty`
	PrePushHook           string `yaml:"pre-push",omitempty`
	PreRebaseHook         string `yaml:"pre-rebase",omitempty`
	PreReceiveHook        string `yaml:"pre-receive",omitempty`
	PrepareCommitMsgHook  string `yaml:"prepare-commit-msg",omitempty`
	UpdateHook            string `yaml:"update",omitempty`
}

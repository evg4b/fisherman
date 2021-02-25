package configuration

type CommitMsgHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type ApplyPatchMsgHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type FsMonitorWatchmanHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PostUpdateHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PreApplyPatchHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PreCommitHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PrePushHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PreRebaseHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type PreReceiveHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

type UpdateHookConfig struct {
	CommonConfig `yaml:"-,inline"`
}

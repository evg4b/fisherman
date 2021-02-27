package configuration

type CommitMsgHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type ApplyPatchMsgHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type FsMonitorWatchmanHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PostUpdateHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PreApplyPatchHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PreCommitHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PrePushHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PreRebaseHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type PreReceiveHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

type UpdateHookConfig struct {
	HookConfig `yaml:"-,inline"`
}

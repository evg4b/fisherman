package hooks

type FsMonitorWatchmanHookConfig struct{}

func (*FsMonitorWatchmanHookConfig) Compile(variables map[string]interface{}) {}

func (*FsMonitorWatchmanHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

func (*FsMonitorWatchmanHookConfig) HasVars() bool {
	return false
}

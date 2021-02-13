package configuration

type FsMonitorWatchmanHookConfig struct{}

func (*FsMonitorWatchmanHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*FsMonitorWatchmanHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

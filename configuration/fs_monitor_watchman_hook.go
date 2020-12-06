package configuration

type FsMonitorWatchmanHookConfig struct{}

func (*FsMonitorWatchmanHookConfig) Compile(variables map[string]interface{}) {}

func (*FsMonitorWatchmanHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

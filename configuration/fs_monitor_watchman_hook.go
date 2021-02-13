package configuration

import "fisherman/internal/expression"

type FsMonitorWatchmanHookConfig struct{}

func (*FsMonitorWatchmanHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*FsMonitorWatchmanHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

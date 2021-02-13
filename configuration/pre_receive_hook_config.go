package configuration

import "fisherman/internal/expression"

type PreReceiveHookConfig struct{}

func (*PreReceiveHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*PreReceiveHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

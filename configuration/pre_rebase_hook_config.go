package configuration

import "fisherman/internal/expression"

type PreRebaseHookConfig struct{}

func (*PreRebaseHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*PreRebaseHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

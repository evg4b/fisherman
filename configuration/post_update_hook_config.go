package configuration

import "fisherman/internal/expression"

type PostUpdateHookConfig struct{}

func (*PostUpdateHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*PostUpdateHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

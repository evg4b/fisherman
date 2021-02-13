package configuration

import "fisherman/internal/expression"

type UpdateHookConfig struct{}

func (*UpdateHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*UpdateHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

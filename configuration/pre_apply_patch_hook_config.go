package configuration

import "fisherman/internal/expression"

type PreApplyPatchHookConfig struct{}

func (*PreApplyPatchHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*PreApplyPatchHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

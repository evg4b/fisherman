package configuration

import "fisherman/internal/expression"

type ApplyPatchMsgHookConfig struct{}

func (*ApplyPatchMsgHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	panic("not supported")
}

func (*ApplyPatchMsgHookConfig) GetVariables() map[string]interface{} {
	panic("not supported")
}

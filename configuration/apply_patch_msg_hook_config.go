package configuration

type ApplyPatchMsgHookConfig struct{}

func (*ApplyPatchMsgHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*ApplyPatchMsgHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

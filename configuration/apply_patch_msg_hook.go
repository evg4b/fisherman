package configuration

type ApplyPatchMsgHookConfig struct{}

func (*ApplyPatchMsgHookConfig) Compile(variables map[string]interface{}) {}

func (*ApplyPatchMsgHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

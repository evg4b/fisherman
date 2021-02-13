package configuration

type PreApplyPatchHookConfig struct{}

func (*PreApplyPatchHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*PreApplyPatchHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

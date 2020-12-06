package configuration

type PreApplyPatchHookConfig struct{}

func (*PreApplyPatchHookConfig) Compile(variables map[string]interface{}) {}

func (*PreApplyPatchHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

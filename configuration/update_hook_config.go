package configuration

type UpdateHookConfig struct{}

func (*UpdateHookConfig) Compile(variables map[string]interface{}) {}

func (*UpdateHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

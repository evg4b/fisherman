package configuration

type UpdateHookConfig struct{}

func (*UpdateHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*UpdateHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

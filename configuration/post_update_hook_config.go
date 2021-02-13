package configuration

type PostUpdateHookConfig struct{}

func (*PostUpdateHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*PostUpdateHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

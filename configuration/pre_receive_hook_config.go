package configuration

type PreReceiveHookConfig struct{}

func (*PreReceiveHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*PreReceiveHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

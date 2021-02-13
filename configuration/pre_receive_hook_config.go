package configuration

type PreReceiveHookConfig struct{}

func (*PreReceiveHookConfig) Compile(variables map[string]interface{}) {}

func (*PreReceiveHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

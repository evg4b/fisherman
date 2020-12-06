package configuration

type PostUpdateHookConfig struct{}

func (*PostUpdateHookConfig) Compile(variables map[string]interface{}) {}

func (*PostUpdateHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

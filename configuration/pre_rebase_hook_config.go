package configuration

type PreRebaseHookConfig struct{}

func (*PreRebaseHookConfig) Compile(variables map[string]interface{}) {}

func (*PreRebaseHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

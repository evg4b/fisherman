package configuration

type PreRebaseHookConfig struct{}

func (*PreRebaseHookConfig) Compile(variables map[string]interface{}) {
	panic("not supported")
}

func (*PreRebaseHookConfig) GetVariablesConfig() VariablesConfig {
	panic("not supported")
}

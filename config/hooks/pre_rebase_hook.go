package hooks

type PreRebaseHookConfig struct{}

func (*PreRebaseHookConfig) Compile(variables map[string]interface{}) {}

func (*PreRebaseHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

func (*PreRebaseHookConfig) HasVars() bool {
	return false
}

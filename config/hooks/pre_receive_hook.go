package hooks

type PreReceiveHookConfig struct{}

func (*PreReceiveHookConfig) Compile(variables map[string]interface{}) {}

func (*PreReceiveHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

func (*PreReceiveHookConfig) HasVars() bool {
	return false
}

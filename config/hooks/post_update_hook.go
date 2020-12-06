package hooks

type PostUpdateHookConfig struct{}

func (*PostUpdateHookConfig) Compile(variables map[string]interface{}) {}

func (*PostUpdateHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

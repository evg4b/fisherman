package hooks

type UpdateHookConfig struct{}

func (*UpdateHookConfig) Compile(variables map[string]interface{}) {}

func (*UpdateHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

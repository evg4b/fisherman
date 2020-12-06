package hooks

type ApplyPatchMsgHookConfig struct{}

func (*ApplyPatchMsgHookConfig) Compile(variables map[string]interface{}) {}

func (*ApplyPatchMsgHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

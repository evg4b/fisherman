package hooks

type PreApplyPatchHookConfig struct{}

func (*PreApplyPatchHookConfig) Compile(variables map[string]interface{}) {}

func (*PreApplyPatchHookConfig) GetVarsSection() Variables {
	panic("not supported")
}

func (*PreApplyPatchHookConfig) HasVars() bool {
	return false
}

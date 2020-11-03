package hooks

// PreCommitHookConfig is structure to storage user configuration for pre-commit hook
type PreCommitHookConfig struct {
	Shell ShellScriptsConfig
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

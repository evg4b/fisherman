package hooks

// PreCommitHookConfig is structure to storage user configuration for pre-commit hook
type PreCommitHookConfig struct {
	Variables Variables          `yaml:"variables,omitempty"`
	Shell     ShellScriptsConfig `yaml:"shell,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

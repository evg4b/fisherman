package hooks

// PreCommitHookConfig is structure to storage user configuration for pre-commit hook
type PreCommitHookConfig struct {
	Variables       Variables     `yaml:"variables,omitempty"`
	Shell           ScriptsConfig `yaml:"shell,omitempty"`
	AddFilesToIndex []string      `yaml:"add-to-index,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

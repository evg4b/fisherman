package hooks

// PrePushHookConfig is structure to storage user configuration for pre-push hook
type PrePushHookConfig struct {
	Variables Variables     `yaml:"variables,omitempty"`
	Shell     ScriptsConfig `yaml:"shell,omitempty"`
}

func (config *PrePushHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

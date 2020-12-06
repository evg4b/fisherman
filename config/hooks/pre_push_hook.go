package hooks

type PrePushHookConfig struct {
	Variables Variables     `yaml:"variables,omitempty"`
	Shell     ScriptsConfig `yaml:"shell,omitempty"`
}

func (config *PrePushHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PrePushHookConfig) GetVarsSection() Variables {
	return config.Variables
}

func (config *PrePushHookConfig) IsEmpty() bool {
	return len(config.Shell) == 0 && config.Variables == Variables{}
}

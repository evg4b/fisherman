package configuration

type PrePushHookConfig struct {
	Variables VariablesConfig `yaml:"variables,omitempty"`
	Shell     ScriptsConfig   `yaml:"shell,omitempty"`
}

func (config *PrePushHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PrePushHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

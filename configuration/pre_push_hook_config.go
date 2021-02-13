package configuration

type PrePushHookConfig struct {
	RulesSection `yaml:"-,inline"`
	Variables    VariablesConfig `yaml:"variables,omitempty"`
}

func (config *PrePushHookConfig) Compile(variables map[string]interface{}) {
}

func (config *PrePushHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

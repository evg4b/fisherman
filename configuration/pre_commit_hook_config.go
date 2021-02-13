package configuration

type PreCommitHookConfig struct {
	RulesSection `yaml:"-,inline"`
	Variables    VariablesConfig `yaml:"variables,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
}

func (config *PreCommitHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

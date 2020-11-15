package hooks

type PreCommitHookConfig struct {
	Variables       Variables     `yaml:"variables,omitempty"`
	Shell           ScriptsConfig `yaml:"shell,omitempty"`
	AddFilesToIndex []string      `yaml:"add-to-index,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PreCommitHookConfig) GetVarsSection() Variables {
	return config.Variables
}

func (*PreCommitHookConfig) HasVars() bool {
	return true
}

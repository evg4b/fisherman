package hooks

import "fisherman/actions"

type PreCommitHookConfig struct {
	Variables       Variables      `yaml:"variables,omitempty"`
	Shell           ScriptsConfig  `yaml:"shell,omitempty"`
	AddFilesToIndex []actions.Glob `yaml:"add-to-index,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PreCommitHookConfig) GetVarsSection() Variables {
	return config.Variables
}

func (config *PreCommitHookConfig) IsEmpty() bool {
	return len(config.Shell) == 0 &&
		len(config.AddFilesToIndex) == 0 &&
		config.Variables == Variables{}
}

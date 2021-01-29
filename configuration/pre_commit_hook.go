package configuration

import "fisherman/actions"

type PreCommitHookConfig struct {
	RulesSection    `yaml:"-,inline"`
	Variables       VariablesConfig `yaml:"variables,omitempty"`
	Shell           ScriptsConfig   `yaml:"shell,omitempty"`
	AddFilesToIndex []actions.Glob  `yaml:"add-to-index,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PreCommitHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

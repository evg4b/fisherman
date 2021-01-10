package configuration

import "fisherman/actions"

type PreCommitHookConfig struct {
	RulesSection        `yaml:"-,inline"`
	Variables           VariablesConfig                    `yaml:"variables,omitempty"`
	Shell               ScriptsConfig                      `yaml:"shell,omitempty"`
	AddFilesToIndex     []actions.Glob                     `yaml:"add-to-index,omitempty"`
	SuppressCommitFiles actions.SuppresCommitFilesSections `yaml:"suppress-commit-files,omitempty"`
}

func (config *PreCommitHookConfig) Compile(variables map[string]interface{}) {
	config.Shell.Compile(variables)
}

func (config *PreCommitHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

func (config *PreCommitHookConfig) IsEmpty() bool {
	return len(config.Shell) == 0 &&
		len(config.AddFilesToIndex) == 0 &&
		config.Variables == VariablesConfig{} &&
		len(config.SuppressCommitFiles.Globs) == 0
}

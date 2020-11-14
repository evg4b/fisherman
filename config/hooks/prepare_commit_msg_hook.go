package hooks

import "fisherman/utils"

// PrepareCommitMsgHookConfig config section for configure prepare-commit-msg hook
type PrepareCommitMsgHookConfig struct {
	Variables Variables `yaml:"variables,omitempty"`
	Message   string    `yaml:"message,omitempty"`
}

func (config *PrepareCommitMsgHookConfig) Compile(variables map[string]interface{}) {
	utils.FillTemplate(&config.Message, variables)
}

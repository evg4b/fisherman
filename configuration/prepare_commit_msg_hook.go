package configuration

import "fisherman/utils"

type PrepareCommitMsgHookConfig struct {
	Variables VariablesConfig `yaml:"variables,omitempty"`
	Message   string          `yaml:"message,omitempty"`
}

func (config *PrepareCommitMsgHookConfig) Compile(variables map[string]interface{}) {
	config.Message = utils.FillTemplate(config.Message, variables)
}

func (config *PrepareCommitMsgHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

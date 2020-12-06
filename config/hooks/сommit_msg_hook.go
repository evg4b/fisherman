package hooks

import "fisherman/utils"

type CommitMsgHookConfig struct {
	Variables     Variables `yaml:"variables,omitempty"`
	NotEmpty      bool      `yaml:"not-empty,omitempty"`
	MessageRegexp string    `yaml:"commit-regexp,omitempty"`
	MessagePrefix string    `yaml:"commit-prefix,omitempty"`
	MessageSuffix string    `yaml:"commit-suffix,omitempty"`
	StaticMessage string    `yaml:"static-message,omitempty"`
}

func (config *CommitMsgHookConfig) Compile(variables map[string]interface{}) {
	utils.FillTemplate(&config.MessagePrefix, variables)
	utils.FillTemplate(&config.MessageSuffix, variables)
	utils.FillTemplate(&config.StaticMessage, variables)
}

func (config *CommitMsgHookConfig) GetVarsSection() Variables {
	return config.Variables
}

func (config *CommitMsgHookConfig) IsEmpty() bool {
	return (*config) == CommitMsgHookConfig{}
}

package configuration

import "fisherman/utils"

type CommitMsgHookConfig struct {
	RulesSection  `yaml:"-,inline"`
	Variables     VariablesConfig `yaml:"variables,omitempty"`
	NotEmpty      bool            `yaml:"not-empty,omitempty"`
	MessageRegexp string          `yaml:"commit-regexp,omitempty"`
	MessagePrefix string          `yaml:"commit-prefix,omitempty"`
	MessageSuffix string          `yaml:"commit-suffix,omitempty"`
	StaticMessage string          `yaml:"static-message,omitempty"`
}

func (config *CommitMsgHookConfig) Compile(vars Variables) {
	config.MessagePrefix = utils.FillTemplate(config.MessagePrefix, vars)
	config.MessageSuffix = utils.FillTemplate(config.MessageSuffix, vars)
	config.StaticMessage = utils.FillTemplate(config.StaticMessage, vars)
}

func (config *CommitMsgHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

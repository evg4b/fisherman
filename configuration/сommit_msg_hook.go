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
	utils.FillTemplate(&config.MessagePrefix, vars)
	utils.FillTemplate(&config.MessageSuffix, vars)
	utils.FillTemplate(&config.StaticMessage, vars)
}

func (config *CommitMsgHookConfig) GetVariablesConfig() VariablesConfig {
	return config.Variables
}

func (config *CommitMsgHookConfig) IsEmpty() bool {
	return !config.NotEmpty &&
		utils.IsEmpty(config.MessageRegexp) &&
		utils.IsEmpty(config.MessagePrefix) &&
		utils.IsEmpty(config.MessageSuffix) &&
		utils.IsEmpty(config.StaticMessage) &&
		len(config.Rules) == 0 &&
		utils.IsEmpty(config.Variables.FromBranch) &&
		utils.IsEmpty(config.Variables.FromLastTag)
}

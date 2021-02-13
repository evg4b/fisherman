package configuration

import (
	"fisherman/internal/expression"
	"fisherman/utils"
)

type CommitMsgHookConfig struct {
	RulesSection     `yaml:"-,inline"`
	VariablesSection `yaml:"-,inline"`
	NotEmpty         bool   `yaml:"not-empty,omitempty"`
	MessageRegexp    string `yaml:"commit-regexp,omitempty"`
	MessagePrefix    string `yaml:"commit-prefix,omitempty"`
	MessageSuffix    string `yaml:"commit-suffix,omitempty"`
	StaticMessage    string `yaml:"static-message,omitempty"`
}

func (config *CommitMsgHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	config.VariablesSection.Compile(engine, global)
	variables := config.VariablesSection.GetVariables()

	config.MessagePrefix = utils.FillTemplate(config.MessagePrefix, global, variables)
	config.MessageSuffix = utils.FillTemplate(config.MessageSuffix, global, variables)
	config.StaticMessage = utils.FillTemplate(config.StaticMessage, global, variables)
}

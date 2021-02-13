package configuration

import (
	"fisherman/internal/expression"
	"fisherman/utils"
)

type PrepareCommitMsgHookConfig struct {
	VariablesSection `yaml:"-,inline"`
	Message          string `yaml:"message,omitempty"`
}

func (config *PrepareCommitMsgHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	config.Message = utils.FillTemplate(config.Message, global)
}

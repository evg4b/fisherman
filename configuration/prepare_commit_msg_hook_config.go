package configuration

import (
	"fisherman/internal/expression"
	"fisherman/utils"
)

type PrepareCommitMsgHookConfig struct {
	VariablesSection `yaml:"-,inline"`
	Message          string `yaml:"message,omitempty"`
}

func (config *PrepareCommitMsgHookConfig) Compile(engine expression.Engine, global map[string]interface{}) error {
	err := config.VariablesSection.Compile(engine, global)
	if err != nil {
		return err
	}

	utils.FillTemplate(&config.Message, config.VariablesSection.GetVariables())

	return nil
}

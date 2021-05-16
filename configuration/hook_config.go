package configuration

import (
	"fisherman/internal/expression"
)

type HookConfig struct {
	VariablesSection `yaml:",inline"`
	RulesSection     `yaml:",inline"`
}

func (config *HookConfig) Compile(engine expression.Engine, global map[string]interface{}) error {
	variables, err := config.VariablesSection.Compile(engine, global)
	if err != nil {
		return err
	}

	config.RulesSection.Compile(variables)

	return nil
}

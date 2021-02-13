package configuration

import "fisherman/internal/expression"

type PreCommitHookConfig struct {
	VariablesSection `yaml:"-,inline"`
	RulesSection     `yaml:"-,inline"`
}

func (config *PreCommitHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {
	// config.VariablesSection.Compile(variables)
}

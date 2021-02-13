package configuration

import "fisherman/internal/expression"

type PrePushHookConfig struct {
	VariablesSection `yaml:"-,inline"`
	RulesSection     `yaml:"-,inline"`
}

func (config *PrePushHookConfig) Compile(engine expression.Engine, global map[string]interface{}) {

}

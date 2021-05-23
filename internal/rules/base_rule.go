package rules

import "fisherman/utils"

// TODO: Create more flexible approach to create correct execution order
var (
	PreScripts  byte = 1
	Scripts     byte = 2
	PostScripts byte = 3
)

type BaseRule struct {
	Type      string `yaml:"type,omitempty"`
	Condition string `yaml:"when,omitempty"`
}

func (rule *BaseRule) GetType() string {
	return rule.Type
}

func (rule *BaseRule) GetContition() string {
	return rule.Condition
}

func (rule *BaseRule) GetPosition() byte {
	return PreScripts
}

func (rule *BaseRule) Compile(variables map[string]interface{}) {
	utils.FillTemplate(&rule.Condition, variables)
}

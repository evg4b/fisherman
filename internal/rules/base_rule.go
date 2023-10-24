package rules

import (
	"fisherman/internal"
	"fisherman/internal/utils"
	"fisherman/internal/validation"

	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
)

var (
	PreScripts  byte = 1
	Scripts     byte = 2
	PostScripts byte = 3
)

type RuleOption = func(rule *BaseRule)

type BaseRule struct {
	Type      string `yaml:"type,omitempty"`
	Condition string `yaml:"when,omitempty"`

	cwd  string
	fs   billy.Filesystem
	repo internal.Repository
	args []string
	env  []string
}

func (rule *BaseRule) GetType() string {
	return rule.Type
}

func (rule *BaseRule) GetPrefix() string {
	return rule.Type
}

func (rule *BaseRule) GetContition() string {
	return rule.Condition
}

func (rule *BaseRule) GetPosition() byte {
	return PreScripts
}

func (rule *BaseRule) Compile(variables map[string]any) {
	utils.FillTemplate(&rule.Condition, variables)
}

func (rule *BaseRule) errorf(message string, a ...any) error {
	return validation.Errorf(rule.GetPrefix(), message, a...)
}

func (rule *BaseRule) Configure(options ...RuleOption) {
	for _, option := range options {
		option(rule)
	}
}

func (rule *BaseRule) arg(index int) (string, error) {
	if index < 0 {
		return "", errors.New("incorrect argument index")
	}

	if rule.args == nil || len(rule.args) <= index {
		return "", errors.Errorf("argument at index %d is not provided", index)
	}

	return rule.args[index], nil
}

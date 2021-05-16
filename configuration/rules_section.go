package configuration

import (
	"fisherman/internal"
	"io"

	"gopkg.in/yaml.v3"
)

// TODO: Add new method in Rule interface to Decode rule from map[string]interface{} and
// try implement comman realization in base rule structure.
type Rule interface {
	GetType() string
	GetContition() string
	GetPosition() byte
	Check(internal.ExecutionContext, io.Writer) error
	Compile(map[string]interface{})
}

type RulesSection struct {
	Rules []Rule
}

type rulesSectionDef struct {
	Rules []ruleDef `yaml:"rules"`
}

func (section *RulesSection) UnmarshalYAML(value *yaml.Node) error {
	section.Rules = []Rule{}
	var definition = rulesSectionDef{}
	err := value.Decode(&definition)
	if err != nil {
		return err
	}

	for _, ruleDef := range definition.Rules {
		section.Rules = append(section.Rules, ruleDef.Rule)
	}

	return nil
}

func (section *RulesSection) Compile(variables map[string]interface{}) {
	for _, rule := range section.Rules {
		rule.Compile(variables)
	}
}

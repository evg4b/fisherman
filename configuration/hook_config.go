package configuration

import (
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/utils"
	"io"

	"github.com/imdario/mergo"
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

type Variables = map[string]interface{}

type HookConfig struct {
	StaticVariables  map[string]string
	ExtractVariables []string
	Rules            []Rule
}

type hookConfigDef struct {
	StaticVariables  map[string]string `yaml:"variables,omitempty"`
	ExtractVariables []string          `yaml:"extract-variables,omitempty"`
	Rules            []ruleDef         `yaml:"rules"`
}

func (c *HookConfig) UnmarshalYAML(value *yaml.Node) error {
	var def hookConfigDef
	err := value.Decode(&def)
	if err != nil {
		return err
	}

	c.ExtractVariables = def.ExtractVariables
	c.StaticVariables = def.StaticVariables

	for _, ruleDef := range def.Rules {
		c.Rules = append(c.Rules, ruleDef.Rule)
	}

	return nil
}

func (c *HookConfig) Compile(engine expression.Engine, global Variables) (Variables, error) {
	variables := map[string]interface{}{}
	err := mergo.MergeWithOverwrite(&variables, global)
	if err != nil {
		return variables, err
	}

	if c.StaticVariables != nil {
		utils.FillTemplatesMap(c.StaticVariables, variables)

		interfaceMap := utils.StringMapToInterfaceMap(c.StaticVariables)
		err = mergo.MergeWithOverwrite(&variables, interfaceMap)
		if err != nil {
			return variables, err
		}
	}

	if c.ExtractVariables != nil {
		utils.FillTemplatesArray(c.ExtractVariables, variables)

		for _, value := range c.ExtractVariables {
			extractedVariables, err := engine.EvalMap(value, variables)
			if err != nil {
				return variables, err
			}

			err = mergo.MergeWithOverwrite(&variables, extractedVariables)
			if err != nil {
				return variables, err
			}
		}
	}

	for _, rule := range c.Rules {
		rule.Compile(variables)
	}

	return variables, nil
}

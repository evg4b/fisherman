package configuration

import (
	"fisherman/internal/expression"
	"fisherman/utils"

	"github.com/imdario/mergo"
)

type Variables = map[string]interface{}

type VariablesSection struct {
	StaticVariables  map[string]string `yaml:"variables,omitempty"`
	ExtractVariables []string          `yaml:"extract-variables,omitempty"`
}

func (config *VariablesSection) Compile(engine expression.Engine, globalVariables Variables) (Variables, error) {
	variables := map[string]interface{}{}
	err := mergo.MergeWithOverwrite(&variables, globalVariables)
	if err != nil {
		return variables, err
	}

	if config.StaticVariables != nil {
		utils.FillTemplatesMap(config.StaticVariables, variables)

		interfaceMap := utils.StringMapToInterfaceMap(config.StaticVariables)
		err = mergo.MergeWithOverwrite(&variables, interfaceMap)
		if err != nil {
			return variables, err
		}
	}

	if config.ExtractVariables != nil {
		utils.FillTemplatesArray(config.ExtractVariables, variables)

		for _, value := range config.ExtractVariables {
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

	return variables, nil
}

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
	compiled         map[string]interface{}
}

func (config *VariablesSection) Compile(engine expression.Engine, globalVariables map[string]interface{}) error {
	config.compiled = map[string]interface{}{}
	err := mergo.MergeWithOverwrite(&config.compiled, globalVariables)
	if err != nil {
		return err
	}

	if config.StaticVariables != nil {
		utils.FillTemplatesMap(config.StaticVariables, config.compiled)

		interfaceMap := utils.StringMapToInterfaceMap(config.StaticVariables)
		err = mergo.MergeWithOverwrite(&config.compiled, interfaceMap)
		if err != nil {
			return err
		}
	}

	if config.ExtractVariables != nil {
		utils.FillTemplatesArray(config.ExtractVariables, config.compiled)

		for _, value := range config.ExtractVariables {
			extractedVariables, err := engine.EvalMap(value, config.compiled)
			if err != nil {
				return err
			}

			err = mergo.MergeWithOverwrite(&config.compiled, extractedVariables)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (config *VariablesSection) GetVariables() map[string]interface{} {
	if config.compiled == nil {
		panic("config id not compiled")
	}

	return config.compiled
}

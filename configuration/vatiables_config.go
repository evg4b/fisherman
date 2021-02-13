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

func (config *VariablesSection) Compile(engine expression.Engine, globalVariables map[string]interface{}) {
	for key, value := range config.StaticVariables {
		config.StaticVariables[key] = utils.FillTemplate(value, globalVariables)
	}

	combinedVariables := map[string]interface{}{}
	err := mergo.MergeWithOverwrite(&combinedVariables, globalVariables)
	if err != nil {
		panic(err)
	}

	if config.StaticVariables != nil {
		interfaceMap := utils.StringMapToInterfaceMap(config.StaticVariables)
		err = mergo.MergeWithOverwrite(&combinedVariables, interfaceMap)
		if err != nil {
			panic(err)
		}
	}

	config.compiled = map[string]interface{}{}
	err = mergo.MergeWithOverwrite(&config.compiled, combinedVariables)
	if err != nil {
		panic(err)
	}

	filledExtractVariables := []string{}
	for _, value := range config.ExtractVariables {
		filledValue := utils.FillTemplate(value, combinedVariables)
		filledExtractVariables = append(filledExtractVariables, filledValue)
		extractedVariables, err := engine.EvalMap(filledValue, combinedVariables)
		if err != nil {
			panic(err)
		}

		err = mergo.MergeWithOverwrite(&config.compiled, extractedVariables)
		if err != nil {
			panic(err)
		}
	}

	config.ExtractVariables = filledExtractVariables
}

func (config *VariablesSection) GetVariables() map[string]interface{} {
	if config.compiled == nil {
		panic("config id not compiled")
	}

	return config.compiled
}

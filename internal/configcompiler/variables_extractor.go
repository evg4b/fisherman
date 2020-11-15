package configcompiler

import (
	"fisherman/config/hooks"
	"fisherman/constants"
	"fisherman/infrastructure"

	"github.com/imdario/mergo"
)

type sourceLoader = func() (string, error)
type variablesLoader = func(string) (map[string]interface{}, error)

type ConfigExtractor struct {
	repository      infrastructure.Repository
	variables       map[string]interface{}
	globalVariables map[string]interface{}
	cwd             string
}

func NewConfigExtractor(
	repository infrastructure.Repository,
	globalVariables map[string]interface{},
	cwd string,
) *ConfigExtractor {
	return &ConfigExtractor{
		repository:      repository,
		cwd:             cwd,
		globalVariables: globalVariables,
	}
}

func (ext *ConfigExtractor) Variables(section hooks.Variables) (map[string]interface{}, error) {
	if ext.variables == nil {
		user, err := ext.repository.GetUser()
		if err != nil {
			return nil, err
		}

		ext.variables = map[string]interface{}{
			constants.FishermanVersionVariable: constants.Version,
			constants.CwdVariable:              ext.cwd,
			constants.UserNameVariable:         user.UserName,
			constants.EmailVariable:            user.Email,
		}

		err = mergo.Map(&ext.variables, ext.globalVariables)
		if err != nil {
			return nil, err
		}
	}

	loaders := []struct {
		source sourceLoader
		load   variablesLoader
	}{
		{ext.repository.GetLastTag, section.GetFromTag},
		{ext.repository.GetCurrentBranch, section.GetFromBranch},
	}

	for _, loader := range loaders {
		err := ext.load(loader.source, loader.load)
		if err != nil {
			return nil, err
		}
	}

	return ext.variables, nil
}

func (ext *ConfigExtractor) load(source sourceLoader, load variablesLoader) error {
	sourceString, err := source()
	if err != nil {
		return err
	}

	additionalValues, err := load(sourceString)
	if err != nil {
		return err
	}

	return mergo.MergeWithOverwrite(&ext.variables, additionalValues)
}

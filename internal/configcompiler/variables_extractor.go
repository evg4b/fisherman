package configcompiler

import (
	"fisherman/config/hooks"
	"fisherman/constants"
	"fisherman/infrastructure"

	"github.com/imdario/mergo"
)

type sourceLoader = func() (string, error)
type variablesLoader = func(string) (map[string]interface{}, error)
type sourceLoaderConfig = struct {
	source sourceLoader
	load   variablesLoader
}

type ConfigExtractor struct {
	repository      infrastructure.Repository
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
	user, err := ext.repository.GetUser()
	if err != nil {
		return nil, err
	}

	variables := map[string]interface{}{
		constants.FishermanVersionVariable: constants.Version,
		constants.CwdVariable:              ext.cwd,
		constants.UserNameVariable:         user.UserName,
		constants.EmailVariable:            user.Email,
	}

	err = mergo.Map(&variables, ext.globalVariables)
	if err != nil {
		return nil, err
	}

	loaders := []sourceLoaderConfig{
		{ext.repository.GetLastTag, section.GetFromTag},
		{ext.repository.GetCurrentBranch, section.GetFromBranch},
	}

	for _, loader := range loaders {
		err := ext.load(variables, loader.source, loader.load)
		if err != nil {
			return nil, err
		}
	}

	return variables, nil
}

func (ext *ConfigExtractor) load(variables map[string]interface{}, source sourceLoader, load variablesLoader) error {
	sourceString, err := source()
	if err != nil {
		return err
	}

	additionalValues, err := load(sourceString)
	if err != nil {
		return err
	}

	return mergo.MergeWithOverwrite(&variables, additionalValues)
}

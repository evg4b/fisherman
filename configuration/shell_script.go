package configuration

import (
	"fisherman/utils"
	"fmt"
	"runtime"
)

type ScriptConfig struct {
	Shell     string            `yaml:"shell,omitempty"`
	Commands  []string          `yaml:"commands,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Dir       string            `yaml:"dir,omitempty"`
	Output    bool              `yaml:"output,omitempty"`
	Condition string            `yaml:"condition,omitempty"`
}

type ScriptsConfig map[string]ScriptConfig

const defaultKey = "default"

func (config *ScriptsConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if *config == nil {
		(*config) = ScriptsConfig{}
	}

	var outputConfig string
	err := unmarshal(&outputConfig)
	if err == nil {
		(*config)[defaultKey] = ScriptConfig{
			Commands: []string{outputConfig},
			Env:      map[string]string{},
		}

		return nil
	}

	err = unmarshalOSRelated(unmarshal, config)
	if err == nil {
		return nil
	}

	var scriptConfig ScriptConfig
	err = unmarshal(&scriptConfig)
	if err == nil {
		(*config)[defaultKey] = scriptConfig

		return nil
	}

	type plain ScriptsConfig

	return unmarshal((*plain)(config))
}

func (config *ScriptConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var outputConfig string
	err := unmarshal(&outputConfig)
	if err == nil {
		config.Commands = []string{outputConfig}
		config.Env = map[string]string{}

		return nil
	}

	var commands []string
	if err := unmarshal(&commands); err == nil {
		config.Commands = commands
		config.Env = map[string]string{}

		return nil
	}

	type plain ScriptConfig
	if err = unmarshal((*plain)(config)); err == nil {
		return nil
	}

	return err
}

func (config *ScriptsConfig) Compile(variables map[string]interface{}) {
	for _, shellScript := range *config {
		for key := range shellScript.Commands {
			utils.FillTemplate(&shellScript.Commands[key], variables)
		}
	}
}

func unmarshalOSRelated(unmarshal func(interface{}) error, config *ScriptsConfig) error {
	var scriptsConfigs map[string]ScriptsConfig

	err := unmarshal(&scriptsConfigs)
	if err == nil {
		data, ok := scriptsConfigs[runtime.GOOS]
		if !ok {
			return fmt.Errorf("script for %s os is not defined", runtime.GOOS)
		}

		(*config) = data
	}

	return err
}

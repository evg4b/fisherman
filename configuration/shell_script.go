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

	var systemConfig struct {
		Windows ScriptsConfig `yaml:"windows,omitempty"`
		Linux   ScriptsConfig `yaml:"linux,omitempty"`
		Darwin  ScriptsConfig `yaml:"darwin,omitempty"`
	}

	if err := unmarshal(&systemConfig); err == nil {
		switch runtime.GOOS {
		case "linux":
			(*config) = systemConfig.Linux
		case "windows":
			(*config) = systemConfig.Windows
		case "darwin":
			(*config) = systemConfig.Darwin
		default:
			panic(fmt.Sprintf("System %s is not supported", runtime.GOOS))
		}

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

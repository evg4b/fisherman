package config

import "fisherman/infrastructure/log"

type FishermanConfig struct {
	GlobalVariables GlobalVariables  `yaml:"variables,omitempty"`
	Hooks           HooksConfig      `yaml:"hooks,omitempty"`
	Output          log.OutputConfig `yaml:"output,omitempty"`
}

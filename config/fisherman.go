package config

import "fisherman/infrastructure/logger"

// FishermanConfig is main structure for unmarshal app configuration
type FishermanConfig struct {
	GlobalVariables GlobalVariables     `yaml:"variables,omitempty"`
	Hooks           HooksConfig         `yaml:"hooks,omitempty"`
	Output          logger.OutputConfig `yaml:"output,omitempty"`
}

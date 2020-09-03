package config

import "fisherman/infrastructure/logger"

type FishermanConfig struct {
	Hooks  HooksConfig         `yaml:"hooks,omitempty"`
	Output logger.OutputConfig `yaml:"output,omitempty"`
}

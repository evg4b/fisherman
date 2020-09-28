package config

import "fisherman/infrastructure/logger"

// FishermanConfig is main structure for unmarshal app configuration
type FishermanConfig struct {
	Hooks  HooksConfig         `yaml:"hooks,omitempty"`
	Output logger.OutputConfig `yaml:"output,omitempty"`
}

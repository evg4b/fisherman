package configuration

import (
	"fisherman/infrastructure/log"
)

type FishermanConfig struct {
	GlobalVariables Variables        `yaml:"variables,omitempty"`
	Hooks           HooksConfig      `yaml:"hooks,omitempty"`
	Output          log.OutputConfig `yaml:"output,omitempty"`
	DefaultShell    string           `yaml:"default-shell,omitempty"`
}

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{},
}

var (
	GlobalMode = "global"
	LocalMode  = "local"
	RepoMode   = "repo"
)

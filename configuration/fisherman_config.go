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

var DefaultConfig = `
hooks:
  commit-msg:
    rules:
      - type: commit-message
        prefix: '[fisherman]'
`

var (
	GlobalMode = "global"
	LocalMode  = "local"
	RepoMode   = "repo"
)

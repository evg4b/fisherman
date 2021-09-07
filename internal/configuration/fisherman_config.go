package configuration

import (
	"fisherman/pkg/log"
)

type FishermanConfig struct {
	GlobalVariables Variables        `yaml:"variables,omitempty"`
	Hooks           HooksConfig      `yaml:"hooks,omitempty"`
	Output          log.OutputConfig `yaml:"output,omitempty"`
	DefaultShell    string           `yaml:"default-shell,omitempty"`
}

var DefaultConfig = `
# yaml-language-server: $schema=https://raw.githubusercontent.com/evg4b/fisherman/develop/json-scheme.json

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

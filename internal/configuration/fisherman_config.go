package configuration

import (
	"fisherman/pkg/log"
)

type FishermanConfig struct {
	GlobalVariables map[string]interface{} `yaml:"variables,omitempty"`
	Hooks           HooksConfig            `yaml:"hooks,omitempty"`
	Output          log.OutputConfig       `yaml:"output,omitempty"`
	DefaultShell    string                 `yaml:"default-shell,omitempty"`
}

var DefaultConfig = `# Documentation {{URL}}

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

var ModeOptions = []string{GlobalMode, RepoMode, LocalMode}

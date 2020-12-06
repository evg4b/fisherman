package configuration

import (
	"fisherman/infrastructure/log"
)

type FishermanConfig struct {
	GlobalVariables Variables  `yaml:"variables,omitempty"`
	Hooks           HooksConfig      `yaml:"hooks,omitempty"`
	Output          log.OutputConfig `yaml:"output,omitempty"`
}

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		CommitMsgHook: CommitMsgHookConfig{
			MessagePrefix: "[fisherman]",
			NotEmpty:      true,
		},
	},
}

var (
	GlobalMode = "global"
	LocalMode  = "local"
	RepoMode   = "repo"
)

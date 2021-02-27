package configuration

import (
	"fisherman/infrastructure/log"
	"fisherman/internal/rules"
)

type FishermanConfig struct {
	GlobalVariables Variables        `yaml:"variables,omitempty"`
	Hooks           HooksConfig      `yaml:"hooks,omitempty"`
	Output          log.OutputConfig `yaml:"output,omitempty"`
	DefaultShell    string           `yaml:"default-shell,omitempty"`
}

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		PreCommitHook: &PreCommitHookConfig{
			HookConfig{
				RulesSection: RulesSection{
					Rules: []Rule{
						&rules.CommitMessage{
							Prefix: "[fisherman]",
						},
					},
				},
			},
		},
	},
}

var (
	GlobalMode = "global"
	LocalMode  = "local"
	RepoMode   = "repo"
)

package config

import (
	"fisherman/config/hooks"
)

// DefaultConfig is default configuration for init command
var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		CommitMsgHook: hooks.CommitMsgHookConfig{
			MessagePrefix: "[fisherman]",
		},
	},
}

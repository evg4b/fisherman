package config

import (
	"fisherman/config/hooks"
)

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		CommitMsgHook: hooks.CommitMsgHookConfig{
			MessagePrefix: "[fisherman]",
			NotEmpty:      true,
		},
	},
}

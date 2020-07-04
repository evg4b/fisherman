package config

import (
	"fisherman/common/rules"
)

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		CommitMsgHook: &rules.CommitMsgHookConfig{
			CommitPrefix: "[fisherman]",
		},
	},
}

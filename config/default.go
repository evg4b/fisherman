package config

var DefaultConfig = FishermanConfig{
	Hooks: HooksConfig{
		PreCommitHook: &struct{}{},
		CommitMsgHook: &struct{}{},
	},
}

package hooks

// PrepareCommitMsgHookConfig config section for configure prepare-commit-msg hook
type PrepareCommitMsgHookConfig struct {
	Message      string `yaml:"message,omitempty"`
	BranchRegExp string `yaml:"branch-regexp,omitempty"`
}

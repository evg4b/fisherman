package hooks

// PrepareCommitMsgHookConfig config section for configure prepare-commit-msg hook
type PrepareCommitMsgHookConfig struct {
	Variables Variables `yaml:"variables,omitempty"`
	Message   string    `yaml:"message,omitempty"`
}

package hooks

// CommitMsgHookConfig is structure to storage user configuration about
type CommitMsgHookConfig struct {
	Variables     Variables `yaml:"variables,omitempty"`
	NotEmpty      bool      `yaml:"not-empty,omitempty"`
	MessageRegexp string    `yaml:"commit-regexp,omitempty"`
	MessagePrefix string    `yaml:"commit-prefix,omitempty"`
	MessageSuffix string    `yaml:"commit-suffix,omitempty"`
	StaticMessage string    `yaml:"static-message,omitempty"`
}

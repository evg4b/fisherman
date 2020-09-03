package hooks

// CommitMsgHookConfig is structure to storage user configuration about
type CommitMsgHookConfig struct {
	NotEmpty      bool   `yaml:"not-empty,omitempty"`
	CommitRegexp  string `yaml:"commit-regexp,omitempty"`
	CommitPrefix  string `yaml:"commit-prefix,omitempty"`
	CommitSuffix  string `yaml:"commit-suffix,omitempty"`
	StaticMessage string `yaml:"static-message,omitempty"`
}

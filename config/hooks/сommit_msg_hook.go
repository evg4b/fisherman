package hooks

type CommitMsgHookConfig struct {
	NotEmpty     bool   `yaml:"not-empty,omitempty"`
	CommitRegexp string `yaml:"commit-regexp,omitempty"`
	CommitPrefix string `yaml:"commit-prefix,omitempty"`
	CommitSuffix string `yaml:"commit-suffix,omitempty"`
}

package rules

const SuppressCommitType = "suppress-commit-files"

type SuppressCommitFiles struct {
	BaseRule        `mapstructure:",squash"`
	Globs           []string `mapstructure:"globs"`
	RemoveFromIndex bool     `mapstructure:"remove-from-index"`
}

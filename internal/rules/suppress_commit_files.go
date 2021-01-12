package rules

import (
	"fisherman/internal"
	"io"
)

const SuppressCommitType = "suppress-commit-files"

type SuppressCommitFiles struct {
	BaseRule        `mapstructure:",squash"`
	Globs           []string `mapstructure:"globs"`
	RemoveFromIndex bool     `mapstructure:"remove-from-index"`
}

func (config SuppressCommitFiles) Check(io.Writer, internal.AsyncContext) error {
	return nil
}

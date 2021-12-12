package rules

import (
	"fisherman/internal"

	"github.com/go-git/go-billy/v5"
)

func WithFileSystem(fileSystem billy.Filesystem) RuleOption {
	return func(ac *BaseRule) {
		ac.fs = fileSystem
	}
}

func WithCwd(cwd string) RuleOption {
	return func(ac *BaseRule) {
		ac.cwd = cwd
	}
}

func WithRepository(repository internal.Repository) RuleOption {
	return func(ac *BaseRule) {
		ac.repo = repository
	}
}

func WithArgs(args []string) RuleOption {
	return func(ac *BaseRule) {
		ac.args = args
	}
}

// WithEnv setups environment variables for BaseRule.
func WithEnv(env []string) RuleOption {
	return func(ac *BaseRule) {
		ac.env = env
	}
}

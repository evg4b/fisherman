package handle

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"

	"github.com/go-git/go-billy/v5"
)

type commandOption = func(*Command)

func WithExpressionEngine(engine expression.Engine) commandOption {
	return func(h *Command) {
		h.engine = engine
	}
}

func WithHooksConfig(config *configuration.HooksConfig) commandOption {
	return func(h *Command) {
		h.config = config
	}
}

func WithGlobalVars(globalVars map[string]interface{}) commandOption {
	return func(h *Command) {
		h.globalVars = globalVars
	}
}

func WithCwd(cwd string) commandOption {
	return func(h *Command) {
		h.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) commandOption {
	return func(h *Command) {
		h.fs = fs
	}
}

func WithRepository(repo internal.Repository) commandOption {
	return func(h *Command) {
		h.repo = repo
	}
}

func WithArgs(args []string) commandOption {
	return func(h *Command) {
		h.args = args
	}
}

func WithEnv(env []string) commandOption {
	return func(h *Command) {
		h.env = env
	}
}

func WithWorkersCount(workersCount uint) commandOption {
	return func(h *Command) {
		h.workersCount = workersCount
	}
}

func WithConfigFiles(configFiles map[string]string) commandOption {
	return func(h *Command) {
		h.configFiles = configFiles
	}
}

package handle

import (
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/expression"
	"io"

	"github.com/go-git/go-billy/v5"
)

type CommandOption = func(*Command)

func WithExpressionEngine(engine expression.Engine) CommandOption {
	return func(h *Command) {
		h.engine = engine
	}
}

func WithHooksConfig(config *configuration.HooksConfig) CommandOption {
	return func(h *Command) {
		h.config = config
	}
}

func WithGlobalVars(globalVars map[string]any) CommandOption {
	return func(h *Command) {
		h.globalVars = globalVars
	}
}

func WithCwd(cwd string) CommandOption {
	return func(h *Command) {
		h.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) CommandOption {
	return func(h *Command) {
		h.fs = fs
	}
}

func WithRepository(repo internal.Repository) CommandOption {
	return func(h *Command) {
		h.repo = repo
	}
}

func WithEnv(env []string) CommandOption {
	return func(h *Command) {
		h.env = env
	}
}

func WithWorkersCount(workersCount uint) CommandOption {
	return func(h *Command) {
		h.workersCount = workersCount
	}
}

func WithConfigFiles(configFiles map[string]string) CommandOption {
	return func(h *Command) {
		h.configFiles = configFiles
	}
}

func WithOutput(output io.Writer) CommandOption {
	return func(h *Command) {
		h.output = output
	}
}

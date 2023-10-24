package handling

import (
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/expression"
	"io"

	"github.com/go-git/go-billy/v5"
)

type HandlerOptions = func(*HookHandler)

func WithExpressionEngine(engine expression.Engine) HandlerOptions {
	return func(h *HookHandler) {
		h.engine = engine
	}
}

func WithHooksConfig(configs *configuration.HooksConfig) HandlerOptions {
	return func(h *HookHandler) {
		h.configs = configs
	}
}

func WithGlobalVars(globalVars map[string]any) HandlerOptions {
	return func(h *HookHandler) {
		h.globalVars = globalVars
	}
}

func WithCwd(cwd string) HandlerOptions {
	return func(h *HookHandler) {
		h.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) HandlerOptions {
	return func(h *HookHandler) {
		h.fs = fs
	}
}

func WithRepository(repo internal.Repository) HandlerOptions {
	return func(h *HookHandler) {
		h.repo = repo
	}
}

func WithArgs(args []string) HandlerOptions {
	return func(h *HookHandler) {
		h.args = args
	}
}

func WithEnv(env []string) HandlerOptions {
	return func(h *HookHandler) {
		h.env = env
	}
}

func WithWorkersCount(workersCount uint) HandlerOptions {
	return func(h *HookHandler) {
		h.workersCount = workersCount
	}
}

func WithOutput(output io.Writer) HandlerOptions {
	return func(h *HookHandler) {
		h.output = output
	}
}

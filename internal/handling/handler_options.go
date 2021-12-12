package handling

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"

	"github.com/go-git/go-billy/v5"
)

type handlerOptions = func(*HookHandler)

func WithExpressionEngine(engine expression.Engine) handlerOptions {
	return func(h *HookHandler) {
		h.engine = engine
	}
}

func WithHooksConfig(configs *configuration.HooksConfig) handlerOptions {
	return func(h *HookHandler) {
		h.configs = configs
	}
}

func WithGlobalVars(globalVars Variables) handlerOptions {
	return func(h *HookHandler) {
		h.globalVars = globalVars
	}
}

func WithCwd(cwd string) handlerOptions {
	return func(h *HookHandler) {
		h.cwd = cwd
	}
}

func WithFileSystem(fs billy.Filesystem) handlerOptions {
	return func(h *HookHandler) {
		h.fs = fs
	}
}

func WithRepository(repo internal.Repository) handlerOptions {
	return func(h *HookHandler) {
		h.repo = repo
	}
}

func WithArgs(args []string) handlerOptions {
	return func(h *HookHandler) {
		h.args = args
	}
}

func WithEnv(env []string) handlerOptions {
	return func(h *HookHandler) {
		h.env = env
	}
}

func WithWorkersCount(workersCount uint) handlerOptions {
	return func(h *HookHandler) {
		h.workersCount = workersCount
	}
}

package appcontext

import (
	"context"
	"fisherman/internal"
	"io"

	"github.com/evg4b/linebyline"
	"github.com/go-git/go-billy/v5"
)

type contextOption = func(*ApplicationContext)

func WithFileSystem(fileSystem billy.Filesystem) contextOption {
	return func(ac *ApplicationContext) {
		ac.fs = fileSystem
	}
}

func WithCwd(cwd string) contextOption {
	return func(ac *ApplicationContext) {
		ac.cwd = cwd
	}
}

func WithShell(shell internal.Shell) contextOption {
	return func(ac *ApplicationContext) {
		ac.shell = shell
	}
}

func WithRepository(repository internal.Repository) contextOption {
	return func(ac *ApplicationContext) {
		ac.repo = repository
	}
}

func WithArgs(args []string) contextOption {
	return func(ac *ApplicationContext) {
		ac.args = args
	}
}

func WithOutput(output io.Writer) contextOption {
	return func(ac *ApplicationContext) {
		ac.output = *linebyline.NewWriterGroup(output)
	}
}

func WithBaseContext(ctx context.Context) contextOption {
	return func(ac *ApplicationContext) {
		ac.baseCtx, ac.cancelBaseCtx = context.WithCancel(ctx)
	}
}

package internal

import (
	"context"
	infra "fisherman/infrastructure"
	"io"
)

type CtxFactory = func(args []string, output io.Writer) *Context

type ExecutionContext interface {
	context.Context
	GlobalVariables() (map[string]interface{}, error)
	Files() infra.FileSystem
	Shell() infra.Shell
	Repository() infra.Repository
	Args() []string
	Output() io.Writer
	Message() (string, error)
	Stop()
}

type AppInfo struct {
	Cwd        string
	Executable string
	Configs    map[string]string
}

func NewCtxFactory(ctx context.Context, fileSystem infra.FileSystem, sysShell infra.Shell, repository infra.Repository) CtxFactory {
	return func(args []string, output io.Writer) *Context {
		return NewInternalContext(ctx, fileSystem, sysShell, repository, args, output)
	}
}

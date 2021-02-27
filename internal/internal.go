package internal

import (
	"context"
	"fisherman/infrastructure"
	"io"
)

type CtxFactory = func(args []string, output io.Writer) *Context

type ExecutionContext interface {
	context.Context
	Files() infrastructure.FileSystem
	Shell() infrastructure.Shell
	Repository() infrastructure.Repository
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

func NewCtxFactory(
	ctx context.Context,
	fileSystem infrastructure.FileSystem,
	sysShell infrastructure.Shell,
	repository infrastructure.Repository,
) CtxFactory {
	return func(args []string, output io.Writer) *Context {
		return NewInternalContext(ctx, fileSystem, sysShell, repository, args, output)
	}
}

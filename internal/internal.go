package internal

import (
	"context"
	"fisherman/infrastructure"
	"io"
)

type CtxFactory = func(args []string, output io.Writer) *Context

type SyncContext interface {
	Files() infrastructure.FileSystem
	Shell() infrastructure.Shell
	Repository() infrastructure.Repository
	Args() []string
	Output() io.Writer
	Message() string
}

type AsyncContext interface {
	SyncContext
	context.Context
	Stop()
}

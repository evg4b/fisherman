package internal

import (
	"context"
	"fisherman/infrastructure"
	"fmt"
	"io"
	"time"
)

type Context struct {
	fileSystem          infrastructure.FileSystem
	shell               infrastructure.Shell
	repository          infrastructure.Repository
	args                []string
	output              io.Writer
	baseContext         context.Context
	cancelBaseContext   context.CancelFunc
	commitmessageLoaded bool
	commitMessage       string
}

func NewInternalContext(
	ctx context.Context,
	fileSystem infrastructure.FileSystem,
	sysShell infrastructure.Shell,
	repository infrastructure.Repository,
	args []string,
	output io.Writer,
) *Context {
	var baseContext, cancel = context.WithCancel(ctx)

	return &Context{
		baseContext:       baseContext,
		cancelBaseContext: cancel,
		fileSystem:        fileSystem,
		shell:             sysShell,
		args:              args,
		output:            output,
		repository:        repository,
	}
}

func (ctx *Context) Files() infrastructure.FileSystem {
	return ctx.fileSystem
}

func (ctx *Context) Shell() infrastructure.Shell {
	return ctx.shell
}

func (ctx *Context) Repository() infrastructure.Repository {
	return ctx.repository
}

func (ctx *Context) Args() []string {
	return ctx.args
}

func (ctx *Context) Output() io.Writer {
	return ctx.output
}

func (ctx *Context) Stop() {
	ctx.cancelBaseContext()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.baseContext.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.baseContext.Done()
}

func (ctx *Context) Err() error {
	return ctx.baseContext.Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.baseContext.Value(key)
}

func (ctx *Context) Message() (string, error) {
	if !ctx.commitmessageLoaded {
		messageFilePath, err := ctx.arg(0)
		if err != nil {
			return "", err
		}

		message, err := ctx.fileSystem.Read(messageFilePath)
		if err != nil {
			return "", err
		}

		ctx.commitMessage = message
		ctx.commitmessageLoaded = true
	}

	return ctx.commitMessage, nil
}

func (ctx *Context) arg(index int) (string, error) {
	if ctx.args == nil || len(ctx.args) <= index {
		return "", fmt.Errorf("argument at index %b is not provided", index)
	}

	return ctx.args[index], nil
}

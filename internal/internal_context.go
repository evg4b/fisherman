package internal

import (
	"context"
	"fisherman/infrastructure"
	"io"
	"time"
)

type InternalContext struct {
	fileSystem          infrastructure.FileSystem
	shell               infrastructure.Shell
	repository          infrastructure.Repository
	args                []string
	output              io.Writer
	baseContext         context.Context
	cancelCaseContext   context.CancelFunc
	commitmessageLoaded bool
	commitMessage       string
}

func NewInternalContext(
	ctx context.Context,
	fileSystem infrastructure.FileSystem,
	shell infrastructure.Shell,
	repository infrastructure.Repository,
	args []string,
	output io.Writer,
) *InternalContext {
	var baseContext, cancel = context.WithCancel(ctx)

	return &InternalContext{
		baseContext:       baseContext,
		cancelCaseContext: cancel,
		fileSystem:        fileSystem,
		shell:             shell,
		args:              args,
		output:            output,
		repository:        repository,
	}
}

func (ctx *InternalContext) Files() infrastructure.FileSystem {
	return ctx.fileSystem
}

func (ctx *InternalContext) Shell() infrastructure.Shell {
	return ctx.shell
}

func (ctx *InternalContext) Repository() infrastructure.Repository {
	return ctx.repository
}

func (ctx *InternalContext) Args() []string {
	return ctx.args
}

func (ctx *InternalContext) Output() io.Writer {
	return ctx.output
}

func (ctx *InternalContext) Stop() {
	ctx.cancelCaseContext()
}

func (ctx *InternalContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.baseContext.Deadline()
}

func (ctx *InternalContext) Done() <-chan struct{} {
	return ctx.baseContext.Done()
}

func (ctx *InternalContext) Err() error {
	return ctx.baseContext.Err()
}

func (ctx *InternalContext) Value(key interface{}) interface{} {
	return ctx.baseContext.Value(key)
}

func (ctx *InternalContext) Message() string {
	if ctx.commitmessageLoaded {
		return ctx.commitMessage
	}

	message, err := ctx.fileSystem.Read(ctx.args[0])
	if err != nil {
		panic(err)
	}

	ctx.commitMessage = message
	ctx.commitmessageLoaded = true

	return message
}

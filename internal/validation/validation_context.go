package validation

import (
	"context"
	"fisherman/infrastructure"
	"io"
	"time"
)

type ValidationContext struct {
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

func NewValidationContext(
	ctx context.Context,
	fileSystem infrastructure.FileSystem,
	shell infrastructure.Shell,
	repository infrastructure.Repository,
	args []string,
	output io.Writer,
) *ValidationContext {
	var baseContext, cancel = context.WithCancel(ctx)

	return &ValidationContext{
		baseContext:       baseContext,
		cancelCaseContext: cancel,
		fileSystem:        fileSystem,
		shell:             shell,
		args:              args,
		output:            output,
		repository:        repository,
	}
}

func (ctx *ValidationContext) Files() infrastructure.FileSystem {
	return ctx.fileSystem
}

func (ctx *ValidationContext) Shell() infrastructure.Shell {
	return ctx.shell
}

func (ctx *ValidationContext) Repository() infrastructure.Repository {
	return ctx.repository
}

func (ctx *ValidationContext) Args() []string {
	return ctx.args
}

func (ctx *ValidationContext) Output() io.Writer {
	return ctx.output
}

func (ctx *ValidationContext) Stop() {
	ctx.cancelCaseContext()
}

func (ctx *ValidationContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.baseContext.Deadline()
}

func (ctx *ValidationContext) Done() <-chan struct{} {
	return ctx.baseContext.Done()
}

func (ctx *ValidationContext) Err() error {
	return ctx.baseContext.Err()
}

func (ctx *ValidationContext) Value(key interface{}) interface{} {
	return ctx.baseContext.Value(key)
}

func (ctx *ValidationContext) Message() string {
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

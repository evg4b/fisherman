package internal

import (
	"context"
	infra "fisherman/infrastructure"
	"fisherman/utils"
	"fmt"
	"io"
	"time"
)

type Context struct {
	fileSystem        infra.FileSystem
	shell             infra.Shell
	repository        infra.Repository
	args              []string
	output            io.Writer
	baseContext       context.Context
	cancelBaseContext context.CancelFunc
}

func NewInternalContext(
	ctx context.Context,
	fileSystem infra.FileSystem,
	sysShell infra.Shell,
	repository infra.Repository,
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

func (ctx *Context) Files() infra.FileSystem {
	return ctx.fileSystem
}

func (ctx *Context) Shell() infra.Shell {
	return ctx.shell
}

func (ctx *Context) Repository() infra.Repository {
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
	messageFilePath, err := ctx.arg(0)
	if err != nil {
		return "", err
	}

	message, err := utils.ReadFileAsString(ctx.fileSystem, messageFilePath)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (ctx *Context) arg(index int) (string, error) {
	if ctx.args == nil || len(ctx.args) <= index {
		return "", fmt.Errorf("argument at index %b is not provided", index)
	}

	return ctx.args[index], nil
}

func (ctx *Context) GlobalVariables() (map[string]interface{}, error) {
	lastTag, err := ctx.repository.GetLastTag()
	if err != nil {
		return nil, err
	}

	currentBranch, err := ctx.repository.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	user, err := ctx.repository.GetUser()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Tag":        lastTag,
		"BranchName": currentBranch,
		"UserEmail":  user.Email,
		"UserName":   user.UserName,
	}, nil
}

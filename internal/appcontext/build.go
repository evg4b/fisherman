package appcontext

import (
	"context"
	i "fisherman/internal"
	"fisherman/pkg/guards"
	"io"
)

type ContextBuilder struct {
	fs     i.FileSystem
	shell  i.Shell
	repo   i.Repository
	args   []string
	output io.Writer
	ctx    context.Context
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		output: io.Discard,
		ctx:    context.TODO(),
		args:   []string{},
	}
}

func (cb *ContextBuilder) WithFileSystem(fileSystem i.FileSystem) *ContextBuilder {
	cb.fs = fileSystem

	return cb
}

func (cb *ContextBuilder) WithShell(shell i.Shell) *ContextBuilder {
	cb.shell = shell

	return cb
}

func (cb *ContextBuilder) WithRepository(repository i.Repository) *ContextBuilder {
	cb.repo = repository

	return cb
}

func (cb *ContextBuilder) WithArgs(args []string) *ContextBuilder {
	cb.args = args

	return cb
}

func (cb *ContextBuilder) WithOutput(output io.Writer) *ContextBuilder {
	cb.output = output

	return cb
}

func (cb *ContextBuilder) WithContext(ctx context.Context) *ContextBuilder {
	cb.ctx = ctx

	return cb
}

func (cb *ContextBuilder) Build() *ApplicationContext {
	guards.ShouldBeDefined(cb.fs, "FileSystem should be connfigured")
	guards.ShouldBeDefined(cb.shell, "Shell should be connfigured")
	guards.ShouldBeDefined(cb.repo, "Repository should be connfigured")

	baseContext, cancelBaseContext := context.WithCancel(cb.ctx)

	return &ApplicationContext{
		fs:            cb.fs,
		shell:         cb.shell,
		repo:          cb.repo,
		args:          cb.args,
		output:        cb.output,
		baseCtx:       baseContext,
		cancelBaseCtx: cancelBaseContext,
	}
}

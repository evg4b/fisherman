package appcontext

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"io"
	"time"

	"github.com/go-errors/errors"
)

type ApplicationContext struct {
	cwd           string
	fs            internal.FileSystem
	shell         internal.Shell
	repo          internal.Repository
	args          []string
	output        io.Writer
	baseCtx       context.Context
	cancelBaseCtx context.CancelFunc
}

const filePathArgumentIndex = 3

func (ctx *ApplicationContext) Files() internal.FileSystem {
	return ctx.fs
}

func (ctx *ApplicationContext) Shell() internal.Shell {
	return ctx.shell
}

func (ctx *ApplicationContext) Repository() internal.Repository {
	return ctx.repo
}

func (ctx *ApplicationContext) Args() []string {
	return ctx.args
}

func (ctx *ApplicationContext) Output() io.Writer {
	return ctx.output
}

func (ctx *ApplicationContext) Stop() {
	ctx.cancelBaseCtx()
}

func (ctx *ApplicationContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.baseCtx.Deadline()
}

func (ctx *ApplicationContext) Done() <-chan struct{} {
	return ctx.baseCtx.Done()
}

func (ctx *ApplicationContext) Err() error {
	return ctx.baseCtx.Err()
}

func (ctx *ApplicationContext) Value(key interface{}) interface{} {
	return ctx.baseCtx.Value(key)
}

func (ctx *ApplicationContext) Message() (string, error) {
	messageFilePath, err := ctx.arg(filePathArgumentIndex)
	if err != nil {
		return "", err
	}

	message, err := utils.ReadFileAsString(ctx.fs, messageFilePath)
	if err != nil {
		return "", err
	}

	return message, nil
}

func (ctx *ApplicationContext) arg(index int) (string, error) {
	if ctx.args == nil || len(ctx.args) <= index {
		return "", errors.Errorf("argument at index %d is not provided", index)
	}

	return ctx.args[index], nil
}

func (ctx *ApplicationContext) GlobalVariables() (map[string]interface{}, error) {
	lastTag, err := ctx.repo.GetLastTag()
	if err != nil {
		return nil, err
	}

	currentBranch, err := ctx.repo.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	user, err := ctx.repo.GetUser()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		constants.UserEmailVariable:        user.Email,
		constants.UserNameVariable:         user.UserName,
		constants.FishermanVersionVariable: constants.Version,
		constants.CwdVariable:              ctx.cwd,
		constants.BranchNameVariable:       currentBranch,
		constants.TagVariable:              lastTag,
	}, nil
}

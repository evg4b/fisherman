package appcontext

import (
	"context"
	"fisherman/internal"
	"fisherman/internal/constants"
	"fisherman/pkg/guards"
	"io"
	"runtime"
	"time"

	"github.com/evg4b/linebyline"
	"github.com/go-errors/errors"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
)

type ApplicationContext struct {
	cwd           string
	fs            billy.Filesystem
	repo          internal.Repository
	args          []string
	env           []string
	output        io.Writer
	baseCtx       context.Context
	cancelBaseCtx context.CancelFunc
}

const filePathArgumentIndex = 3

func NewContext(options ...contextOption) *ApplicationContext {
	baseCtx, cancelBaseCtx := context.WithCancel(context.TODO())
	context := ApplicationContext{
		cwd:           "",
		args:          []string{},
		env:           []string{},
		output:        io.Discard,
		baseCtx:       baseCtx,
		cancelBaseCtx: cancelBaseCtx,
	}

	for _, option := range options {
		option(&context)
	}

	guards.ShouldBeDefined(context.fs, "FileSystem should be connfigured")
	guards.ShouldBeDefined(context.repo, "Repository should be connfigured")

	return &context
}

func (ctx *ApplicationContext) Files() billy.Filesystem {
	return ctx.fs
}

func (ctx *ApplicationContext) Repository() internal.Repository {
	return ctx.repo
}

func (ctx *ApplicationContext) Args() []string {
	return ctx.args
}

func (ctx *ApplicationContext) Output() io.WriteCloser {
	return linebyline.NewByLineWriter(ctx.output)
}

func (ctx *ApplicationContext) Cancel() {
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
	messageFilePath, err := ctx.Arg(filePathArgumentIndex)
	if err != nil {
		return "", err
	}

	message, err := util.ReadFile(ctx.fs, messageFilePath)
	if err != nil {
		return "", err
	}

	return string(message), nil
}

// Cwd returns current working directory.
func (ctx *ApplicationContext) Cwd() string {
	return ctx.cwd
}

func (ctx *ApplicationContext) Arg(index int) (string, error) {
	if index < 0 {
		return "", errors.New("incorrect argument index")
	}

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
		constants.OsVariable:               runtime.GOOS,
	}, nil
}

// Env returns list of environment variables.
func (ctx *ApplicationContext) Env() []string {
	return ctx.env
}

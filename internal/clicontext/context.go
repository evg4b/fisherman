package clicontext

import (
	"context"
	"fisherman/config"
	"fisherman/infrastructure"
	"os/user"
	"time"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config          *config.HooksConfig
	User            *user.User
	App             *AppInfo
	Files           infrastructure.FileSystem
	Repository      infrastructure.Repository
	Shell           infrastructure.Shell
	variables       map[string]interface{}
	globalVariables map[string]interface{}
	base            context.Context
	cancel          context.CancelFunc
}

// AppInfo is application info structure
type AppInfo struct {
	Cwd              string
	Executable       string
	GlobalConfigPath string
	LocalConfigPath  string
	RepoConfigPath   string
}

// Args is structure for params in cli command context constructor
type Args struct {
	FileSystem      infrastructure.FileSystem
	User            *user.User
	App             *AppInfo
	Config          *config.FishermanConfig
	Repository      infrastructure.Repository
	GlobalVariables map[string]interface{}
	Shell           infrastructure.Shell
}

// NewContext constructor for cli command context
func NewContext(baseCtx context.Context, args Args) *CommandContext {
	ctx, cancel := context.WithCancel(baseCtx)

	return &CommandContext{
		Config:          &args.Config.Hooks,
		User:            args.User,
		App:             args.App,
		Files:           args.FileSystem,
		Repository:      args.Repository,
		Shell:           args.Shell,
		globalVariables: args.GlobalVariables,
		base:            ctx,
		cancel:          cancel,
	}
}

func (ctx *CommandContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.base.Deadline()
}

func (ctx *CommandContext) Done() <-chan struct{} {
	return ctx.base.Done()
}

func (ctx *CommandContext) Err() error {
	return ctx.base.Err()
}

func (ctx *CommandContext) Value(key interface{}) interface{} {
	return ctx.base.Value(key)
}

func (ctx *CommandContext) Stop() {
	ctx.cancel()
}

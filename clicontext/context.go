package clicontext

import (
	"fisherman/config"
	"fisherman/infrastructure"
	"os/user"
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
func NewContext(args Args) *CommandContext {
	return &CommandContext{
		Config:          &args.Config.Hooks,
		User:            args.User,
		App:             args.App,
		Files:           args.FileSystem,
		Repository:      args.Repository,
		Shell:           args.Shell,
		globalVariables: args.GlobalVariables,
	}
}

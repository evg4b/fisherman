package clicontext

import (
	"fisherman/config"
	"fisherman/infrastructure"
	"os/user"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config     *config.HooksConfig
	User       *user.User
	App        *AppInfo
	Files      infrastructure.FileAccessor
	Repository infrastructure.Repository
	Shell      infrastructure.Shell
	Variables  map[string]interface{}
}

// AppInfo is application info structure
type AppInfo struct {
	Cwd                string
	Executable         string
	GlobalConfigPath   string
	LocalConfigPath    string
	RepoConfigPath     string
	IsRegisteredInPath bool
}

// Args is structure for params in cli command context constructor
type Args struct {
	FileAccessor infrastructure.FileAccessor
	User         *user.User
	App          *AppInfo
	Config       *config.FishermanConfig
	Repository   infrastructure.Repository
	Variables    map[string]interface{}
	Shell        infrastructure.Shell
}

// NewContext constructor for cli command context
func NewContext(args Args) *CommandContext {
	return &CommandContext{
		&args.Config.Hooks,
		args.User,
		args.App,
		args.FileAccessor,
		args.Repository,
		args.Shell,
		args.Variables,
	}
}

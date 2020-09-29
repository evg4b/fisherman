package commands

import (
	"fisherman/config"
	"fisherman/infrastructure"
	"os/user"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config    *config.HooksConfig
	User      *user.User
	App       *AppInfo
	Files     infrastructure.FileAccessor
	Variables map[string]string
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

// CliCommandContextParams is structure for params in cli command context constructor
type CliCommandContextParams struct {
	FileAccessor infrastructure.FileAccessor
	Usr          *user.User
	App          *AppInfo
	Config       *config.FishermanConfig
	Variables    map[string]string
}

// NewContext constructor for cli command context
func NewContext(args CliCommandContextParams) *CommandContext {
	return &CommandContext{
		&args.Config.Hooks,
		args.Usr,
		args.App,
		args.FileAccessor,
		args.Variables,
	}
}

package commands

import (
	"fisherman/config"
	"fisherman/infrastructure/io"
	"os/user"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config *config.FishermanConfig
	User   *user.User
	App    *AppInfo
	Files  io.FileAccessor
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
	FileAccessor io.FileAccessor
	Usr          *user.User
	App          *AppInfo
	Config       *config.FishermanConfig
}

// NewContext constructor for cli command context
func NewContext(args CliCommandContextParams) *CommandContext {
	return &CommandContext{
		args.Config,
		args.Usr,
		args.App,
		args.FileAccessor,
	}
}

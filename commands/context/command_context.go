package context

import (
	"fisherman/config"
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"fisherman/utils"
	"os/user"
)

// CommandContext is cli context structure
type CommandContext struct {
	Config  *config.FishermanConfig
	User    *user.User
	Files   io.FileAccessor
	Logger  logger.Logger
	AppInfo AppInfo
}

// CliCommandContextParams is structure for params in cli command context constructor
type CliCommandContextParams struct {
	FileAccessor io.FileAccessor
	Usr          *user.User
	Cwd          string
	AppPath      string
	Config       *config.FishermanConfig
	ConfigInfo   *config.LoadInfo
	Logger       logger.Logger
}

// NewContext constructor for cli command context
func NewContext(args CliCommandContextParams) *CommandContext {
	return &CommandContext{
		args.Config,
		args.Usr,
		args.FileAccessor,
		args.Logger,
		AppInfo{
			AppPath:            args.AppPath,
			GlobalConfigPath:   args.ConfigInfo.GlobalConfigPath,
			LocalConfigPath:    args.ConfigInfo.LocalConfigPath,
			IsRegisteredInPath: utils.IsCommandExists(constants.AppConfigName),
			Cwd:                args.Cwd,
		},
	}
}

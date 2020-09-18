package runner

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	commandList  []commands.CliCommand
	fileAccessor io.FileAccessor
	systemUser   *user.User
	logger       logger.Logger
	config       *config.FishermanConfig
	configInfo   *config.LoadInfo
}

// CreateRunnerArgs is structure to pass arguments in constructor
type CreateRunnerArgs struct {
	CommandList  []commands.CliCommand
	FileAccessor io.FileAccessor
	SystemUser   *user.User
	Logger       logger.Logger
	Config       *config.FishermanConfig
	ConfigInfo   *config.LoadInfo
}

// NewRunner is constructor for Runner
func NewRunner(args CreateRunnerArgs) *Runner {

	return &Runner{
		args.CommandList,
		args.FileAccessor,
		args.SystemUser,
		args.Logger,
		args.Config,
		args.ConfigInfo,
	}
}

package runner

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	initc "fisherman/commands/init"
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"flag"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	fileAccessor io.FileAccessor
	systemUser   *user.User
	commandList  []commands.CliCommand
	version      string
	logger       logger.Logger
}

// NewRunner is constructor for Runner
func NewRunner(fileAccessor io.FileAccessor, systemUser *user.User, logger logger.Logger) *Runner {
	version := constants.Version
	return &Runner{
		fileAccessor,
		systemUser,
		[]commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, fileAccessor),
		},
		version,
		logger,
	}
}

package runner

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	initc "fisherman/commands/init"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/logger"
	"flag"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	commandList  []commands.CliCommand
	fileAccessor io.FileAccessor
	systemUser   *user.User
	logger       logger.Logger
}

// NewRunner is constructor for Runner
func NewRunner(fileAccessor io.FileAccessor, systemUser *user.User, logger logger.Logger) *Runner {
	return &Runner{
		[]commands.CliCommand{
			initc.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, fileAccessor),
		},
		fileAccessor,
		systemUser,
		logger,
	}
}

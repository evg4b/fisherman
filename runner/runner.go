package runner

import (
	"fisherman/commands"
	"fisherman/commands/handle"
	"fisherman/commands/init"
	"fisherman/constants"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/reporter"
	"flag"
	"os/user"
)

// Runner is main app structure
type Runner struct {
	fileAccessor io.FileAccessor
	systemUser   *user.User
	reporter     reporter.Reporter
	commandList  []commands.CliCommand
	version      string
}

// NewRunner is constructor for Runner
func NewRunner(fileAccessor io.FileAccessor, systemUser *user.User, reporter reporter.Reporter) *Runner {
	version := constants.Version
	return &Runner{
		fileAccessor,
		systemUser,
		reporter,
		[]commands.CliCommand{
			init.NewCommand(flag.ExitOnError),
			handle.NewCommand(flag.ExitOnError, reporter, fileAccessor),
		},
		version,
	}
}

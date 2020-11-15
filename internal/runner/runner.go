package runner

import (
	"fisherman/commands"
	"fisherman/internal"
)

type Runner struct {
	commands []commands.CliCommand
	app      *internal.AppInfo
}

func NewRunner(commands []commands.CliCommand, app *internal.AppInfo) *Runner {
	return &Runner{commands, app}
}

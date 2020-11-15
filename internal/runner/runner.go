package runner

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/internal"
)

type Runner struct {
	commands []commands.CliCommand
	config   *config.FishermanConfig
	app      *internal.AppInfo
}

func NewRunner(commands []commands.CliCommand, config *config.FishermanConfig, app *internal.AppInfo) *Runner {
	return &Runner{commands, config, app}
}

package runner

import (
	"fisherman/commands"
)

type Runner struct {
	commands []commands.CliCommand
}

func NewRunner(commands []commands.CliCommand) *Runner {
	return &Runner{commands}
}

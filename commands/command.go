package commands

import "fisherman/clicontext"

// CliCommand is base command interface
type CliCommand interface {
	Init(args []string) error
	Run(ctx *clicontext.CommandContext) error
	Name() string
	Description() string
}

package commands

import "fisherman/commands/context"

// CliCommand is base command interface
type CliCommand interface {
	Run(ctx *context.CommandContext, args []string) error
	Name() string
}

package commands

import "fisherman/commands/context"

// CliCommand is base command interface
type CliCommand interface {
	Run(ctx context.Context, args []string) error
	Name() string
}

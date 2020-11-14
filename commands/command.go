package commands

import "fisherman/internal/clicontext"

type CliCommand interface {
	Init(args []string) error
	Run(ctx *clicontext.CommandContext) error
	Name() string
	Description() string
}

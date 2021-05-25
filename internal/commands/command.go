package commands

import i "fisherman/internal"

type CliCommand interface {
	Init(args []string) error
	Run(ctx i.ExecutionContext) error
	Name() string
	Description() string
}

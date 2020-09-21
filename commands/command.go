package commands

// CliCommand is base command interface
type CliCommand interface {
	Init(args []string) error
	Run(ctx *CommandContext) error
	Name() string
}

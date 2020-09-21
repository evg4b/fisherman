package commands

// CliCommand is base command interface
type CliCommand interface {
	Run(ctx *CommandContext, args []string) error
	Name() string
}

package version

import (
	"fisherman/commands"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"flag"
	"fmt"
)

// Command is structure for storage information about remove command
type Command struct {
	fs *flag.FlagSet
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Vertion command created")

	return &Command{
		fs: flag.NewFlagSet("version", handling),
	}
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.fs.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *commands.CommandContext) error {
	_, err := fmt.Fprintln(log.Stdout(), constants.Version)

	return err
}

// Name returns namand name
func (c *Command) Name() string {
	return c.fs.Name()
}

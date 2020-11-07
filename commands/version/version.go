package version

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"flag"
	"fmt"
)

// Command is structure for storage information about remove command
type Command struct {
	flagSet *flag.FlagSet
}

// NewCommand is constructor for init command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Version command created")

	return &Command{
		flagSet: flag.NewFlagSet("version", handling),
	}
}

// Init initialize handle command
func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

// Run executes init command
func (c *Command) Run(ctx *clicontext.CommandContext) error {
	_, err := fmt.Fprintln(log.Stdout(), constants.Version)

	return err
}

// Name returns command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

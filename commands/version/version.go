package version

import (
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal/clicontext"
	"flag"
	"fmt"
)

type Command struct {
	flagSet *flag.FlagSet
	usage   string
}

func NewCommand() *Command {
	return &Command{
		flagSet: flag.NewFlagSet("version", flag.ExitOnError),
		usage:   "prints fisherman version",
	}
}

func (c *Command) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Command) Run(ctx *clicontext.CommandContext) error {
	_, err := fmt.Fprintln(log.Stdout(), constants.Version)

	return err
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

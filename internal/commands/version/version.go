package version

import (
	"context"
	"fisherman/internal/constants"
	"fisherman/pkg/log"
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

func (c *Command) Run(ctx context.Context, _ []string) error {
	_, err := fmt.Fprintf(log.Stdout(), "%s@%s", constants.AppName, constants.Version)

	return err
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

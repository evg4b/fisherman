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

func (command *Command) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *Command) Run(ctx context.Context) error {
	_, err := fmt.Fprintf(log.Stdout(), "%s@%s", constants.AppName, constants.Version)

	return err
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}

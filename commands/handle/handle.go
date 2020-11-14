package handle

import (
	"fisherman/internal/handling"
	"flag"
)

type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers map[string]handling.Handler
	usage    string
}

func NewCommand(handlers map[string]handling.Handler) *Command {
	flagSet := flag.NewFlagSet("handle", flag.ExitOnError)
	command := &Command{
		flagSet:  flagSet,
		handlers: handlers,
		usage:    "starts hook processing based on the config file (for debugging only)",
	}
	flagSet.StringVar(&command.hook, "hook", "<empty>", "hook name")

	return command
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

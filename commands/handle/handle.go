package handle

import (
	"fisherman/config"
	"fisherman/internal"
	"fisherman/internal/handling"
	"flag"
)

type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers map[string]handling.Handler
	usage    string
	config   *config.HooksConfig
	app      *internal.AppInfo
}

func NewCommand(handlers map[string]handling.Handler, config *config.HooksConfig, app *internal.AppInfo) *Command {
	flagSet := flag.NewFlagSet("handle", flag.ExitOnError)
	command := &Command{
		flagSet:  flagSet,
		handlers: handlers,
		usage:    "starts hook processing based on the config file (for debugging only)",
		config:   config,
		app:      app,
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

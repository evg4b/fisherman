package handle

import (
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/hookfactory"
	"flag"
)

type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers hookfactory.HandlerList
	usage    string
	config   *configuration.HooksConfig
	app      *internal.AppInfo
}

func NewCommand(handlers hookfactory.HandlerList, config *configuration.HooksConfig, app *internal.AppInfo) *Command {
	command := &Command{
		flagSet:  flag.NewFlagSet("handle", flag.ExitOnError),
		handlers: handlers,
		usage:    "starts hook processing based on the config file (for debugging only)",
		config:   config,
		app:      app,
	}

	command.flagSet.StringVar(&command.hook, "hook", "<empty>", "hook name")

	return command
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}

package handle

import (
	"fisherman/commands/handle/hooks"
	"fisherman/config"
	"fisherman/internal"
	"flag"
)

type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers hooks.HandlerList
	usage    string
	config   *config.HooksConfig
	app      *internal.AppInfo
}

func NewCommand(handlers hooks.HandlerList, config *config.HooksConfig, app *internal.AppInfo) *Command {
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

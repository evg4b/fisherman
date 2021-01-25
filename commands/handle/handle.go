package handle

import (
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/hookfactory"
	"flag"
)

type Command struct {
	ctxFactory  internal.CtxFactory
	flagSet     *flag.FlagSet
	hook        string
	hookFactory hookfactory.Factory
	usage       string
	config      *configuration.HooksConfig
	app         *internal.AppInfo
}

func NewCommand(
	hookFactory hookfactory.Factory,
	ctxFactory internal.CtxFactory,
	config *configuration.HooksConfig,
	app *internal.AppInfo,
) *Command {
	command := &Command{
		flagSet:     flag.NewFlagSet("handle", flag.ExitOnError),
		hookFactory: hookFactory,
		usage:       "starts hook processing based on the config file (for debugging only)",
		config:      config,
		app:         app,
		ctxFactory:  ctxFactory,
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

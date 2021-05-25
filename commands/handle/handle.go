package handle

import (
	cnfg "fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/handling"
	"flag"
)

type Command struct {
	flagSet     *flag.FlagSet
	hook        string
	hookFactory handling.Factory
	usage       string
	config      *cnfg.HooksConfig
	app         internal.AppInfo
}

func NewCommand(hookFactory handling.Factory, config *cnfg.HooksConfig, app internal.AppInfo) *Command {
	command := &Command{
		flagSet:     flag.NewFlagSet("handle", flag.ExitOnError),
		hookFactory: hookFactory,
		usage:       "starts hook processing based on the config file (for debugging only)",
		config:      config,
		app:         app,
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

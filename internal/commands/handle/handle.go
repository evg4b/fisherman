package handle

import (
	"fisherman/internal"
	cnfg "fisherman/internal/configuration"
	"fisherman/internal/expression"
	"flag"

	"github.com/go-git/go-billy/v5"
)

type Command struct {
	flagSet      *flag.FlagSet
	hook         string
	usage        string
	engine       expression.Engine
	config       *cnfg.HooksConfig
	globalVars   cnfg.Variables
	cwd          string
	fs           billy.Filesystem
	repo         internal.Repository
	args         []string
	env          []string
	workersCount uint
	configFiles  map[string]string
}

func NewCommand(options ...commandOption) *Command {
	command := &Command{
		flagSet: flag.NewFlagSet("handle", flag.ExitOnError),
		usage:   "starts hook processing based on the config file (for debugging only)",
	}
	command.flagSet.StringVar(&command.hook, "hook", "<empty>", "hook name")

	for _, option := range options {
		option(command)
	}

	return command
}

func (command *Command) Name() string {
	return command.flagSet.Name()
}

func (command *Command) Description() string {
	return command.usage
}

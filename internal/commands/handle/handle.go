package handle

import (
	"fisherman/internal"
	"fisherman/internal/configuration"
	"fisherman/internal/expression"
	"flag"
	"io"

	"github.com/go-git/go-billy/v5"
)

type Command struct {
	flagSet      *flag.FlagSet
	hook         string
	usage        string
	engine       expression.Engine
	config       *configuration.HooksConfig
	globalVars   map[string]interface{}
	cwd          string
	fs           billy.Filesystem
	repo         internal.Repository
	args         []string
	env          []string
	workersCount uint
	configFiles  map[string]string
	output       io.Writer
}

const defaultWorkerCount = 5

func NewCommand(options ...commandOption) *Command {
	command := &Command{
		flagSet:      flag.NewFlagSet("handle", flag.ExitOnError),
		usage:        "starts hook processing based on the config file (for debugging only)",
		workersCount: defaultWorkerCount,
		output:       io.Discard,
	}
	command.flagSet.StringVar(&command.hook, "hook", "<empty>", "hook name")

	for _, option := range options {
		option(command)
	}

	return command
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

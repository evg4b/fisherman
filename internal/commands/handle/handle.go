package handle

import (
	"flag"
	"github.com/evg4b/fisherman/internal"
	"github.com/evg4b/fisherman/internal/configuration"
	"github.com/evg4b/fisherman/internal/expression"
	"github.com/evg4b/fisherman/pkg/guards"
	"io"

	"github.com/go-git/go-billy/v5"
)

type Command struct {
	flagSet      *flag.FlagSet
	hook         string
	usage        string
	engine       expression.Engine
	config       *configuration.HooksConfig
	globalVars   map[string]any
	cwd          string
	fs           billy.Filesystem
	repo         internal.Repository
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
		configFiles:  map[string]string{},
		globalVars:   map[string]any{},
		env:          []string{},
	}

	for _, option := range options {
		option(command)
	}

	guards.ShouldBeDefined(command.fs, "FileSystem should be configured")
	guards.ShouldBeDefined(command.repo, "Repository should be configured")
	guards.ShouldBeNotEmpty(command.cwd, "Cwd should be configured")
	guards.ShouldBeDefined(command.engine, "ExpressionEngine should be configured")
	guards.ShouldBeDefined(command.config, "HooksConfig should be configured")

	command.flagSet.StringVar(&command.hook, "hook", "<empty>", "hook name")

	return command
}

func (c *Command) Name() string {
	return c.flagSet.Name()
}

func (c *Command) Description() string {
	return c.usage
}

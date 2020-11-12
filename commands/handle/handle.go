package handle

import (
	"fisherman/handlers"
	"fisherman/infrastructure/log"
	"flag"
)

// Command is structure for storage information about handle command
type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers map[string]handlers.Handler
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling, handlers map[string]handlers.Handler) *Command {
	defer log.Debug("Handle command created")
	flagSet := flag.NewFlagSet("handle", handling)
	command := &Command{
		flagSet:  flagSet,
		handlers: handlers,
	}
	flagSet.StringVar(&command.hook, "hook", "", "")

	return command
}

// Name returns handler command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

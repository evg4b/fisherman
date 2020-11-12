package handle

import (
	"fisherman/clicontext"
	"fisherman/config"
	c "fisherman/constants"
	"fisherman/handlers"
	"fisherman/infrastructure/log"
	"flag"
)

type Handler interface {
	IsConfigured(config *config.HooksConfig) bool
	Handle(ctx *clicontext.CommandContext, args []string) error
}

// Command is structure for storage information about handle command
type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers map[string]Handler
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Handle command created")
	flagSet := flag.NewFlagSet("handle", handling)
	command := &Command{
		flagSet: flagSet,
		handlers: map[string]Handler{
			c.ApplyPatchMsgHook:     new(handlers.NotSupportedHandler),
			c.CommitMsgHook:         new(handlers.CommitMsgHandler),
			c.FsMonitorWatchmanHook: new(handlers.NotSupportedHandler),
			c.PostUpdateHook:        new(handlers.NotSupportedHandler),
			c.PreApplyPatchHook:     new(handlers.NotSupportedHandler),
			c.PreCommitHook:         new(handlers.NotSupportedHandler),
			c.PrePushHook:           new(handlers.PrePushHandler),
			c.PreRebaseHook:         new(handlers.NotSupportedHandler),
			c.PreReceiveHook:        new(handlers.NotSupportedHandler),
			c.PrepareCommitMsgHook:  new(handlers.PrepareCommitMsgHandler),
			c.UpdateHook:            new(handlers.NotSupportedHandler),
		},
	}
	flagSet.StringVar(&command.hook, "hook", "", "")

	return command
}

// Name returns handler command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

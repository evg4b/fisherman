package handle

import (
	c "fisherman/constants"
	"fisherman/handlers"
	"fisherman/infrastructure/log"
	"flag"
)

// Command is structure for storage information about handle command
type Command struct {
	flagSet  *flag.FlagSet
	hook     string
	handlers map[string]handlers.HookHandler
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling) *Command {
	defer log.Debug("Handle command created")
	flagSet := flag.NewFlagSet("handle", handling)
	command := &Command{
		flagSet: flagSet,
		handlers: map[string]handlers.HookHandler{
			c.ApplyPatchMsgHook:     handlers.ApplyPatchMsgHandler,
			c.CommitMsgHook:         handlers.CommitMsgHandler,
			c.FsMonitorWatchmanHook: handlers.FsMonitorWatchmanHandler,
			c.PostUpdateHook:        handlers.PostUpdateHandler,
			c.PreApplyPatchHook:     handlers.PreApplyPatchHandler,
			c.PreCommitHook:         handlers.PreCommitHandler,
			c.PrePushHook:           handlers.PrePushHandler,
			c.PreRebaseHook:         handlers.PreRebaseHandler,
			c.PreReceiveHook:        handlers.PreReceiveHandler,
			c.PrepareCommitMsgHook:  handlers.PrepareCommitMsgHandler,
			c.UpdateHook:            handlers.UpdateHandler,
		},
	}
	flagSet.StringVar(&command.hook, "hook", "", "")

	return command
}

// Name returns command name
func (c *Command) Name() string {
	return c.flagSet.Name()
}

package handle

import (
	c "fisherman/constants"
	"fisherman/handlers"
	"flag"
)

// Command is structure for storage information about handle command
type Command struct {
	fs       *flag.FlagSet
	hook     string
	handlers map[string]handlers.HookHandler
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{
		fs: fs,
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
	fs.StringVar(&c.hook, "hook", "", "")

	return c
}

// Name returns command name
func (c *Command) Name() string {
	return c.fs.Name()
}

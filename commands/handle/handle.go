package handle

import (
	c "fisherman/constants"
	handler "fisherman/handlers"
	"fisherman/handlers/applypatchmsg"
	"fisherman/handlers/commitmsg"
	"fisherman/handlers/fsmonitorwatchman"
	"fisherman/handlers/postupdate"
	"fisherman/handlers/preapplypatch"
	"fisherman/handlers/precommit"
	"fisherman/handlers/preparecommitmsg"
	"fisherman/handlers/prepush"
	"fisherman/handlers/prerebase"
	"fisherman/handlers/prereceive"
	"fisherman/handlers/update"
	"fisherman/infrastructure/io"
	"fisherman/infrastructure/reporter"
	"flag"
)

// Command is structure for storage information about handle command
type Command struct {
	fs       *flag.FlagSet
	hook     string
	args     []string
	reporter reporter.Reporter
	handlers map[string]handler.HookHandler
}

// NewCommand is constructor for handle command
func NewCommand(handling flag.ErrorHandling, r reporter.Reporter, f io.FileAccessor) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{
		fs:       fs,
		reporter: r,
		handlers: map[string]handler.HookHandler{
			c.ApplyPatchMsgHook:     applypatchmsg.NewHandler(),
			c.CommitMsgHook:         commitmsg.NewHandler(),
			c.FsMonitorWatchmanHook: fsmonitorwatchman.NewHandler(),
			c.PostUpdateHook:        postupdate.NewHandler(),
			c.PreApplyPatchHook:     preapplypatch.NewHandler(),
			c.PreCommitHook:         precommit.NewHandler(),
			c.PrePushHook:           prepush.NewHandler(),
			c.PreRebaseHook:         prerebase.NewHandler(),
			c.PreReceiveHook:        prereceive.NewHandler(),
			c.PrepareCommitMsgHook:  preparecommitmsg.NewHandler(f),
			c.UpdateHook:            update.NewHandler(),
		},
	}
	fs.StringVar(&c.hook, "hook", "", "")
	return c
}

// Name returns command name
func (c *Command) Name() string {
	return c.fs.Name()
}

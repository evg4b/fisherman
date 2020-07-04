package handle

import (
	"fisherman/constants"
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
func NewCommand(handling flag.ErrorHandling, reporter reporter.Reporter) *Command {
	fs := flag.NewFlagSet("handle", handling)
	c := &Command{
		fs:       fs,
		reporter: reporter,
		handlers: map[string]handler.HookHandler{
			constants.ApplyPatchMsgHook:     applypatchmsg.Handler,
			constants.CommitMsgHook:         commitmsg.Handler,
			constants.FsMonitorWatchmanHook: fsmonitorwatchman.Handler,
			constants.PostUpdateHook:        postupdate.Handler,
			constants.PreApplyPatchHook:     preapplypatch.Handler,
			constants.PreCommitHook:         precommit.Handler,
			constants.PrePushHook:           prepush.Handler,
			constants.PreRebaseHook:         prerebase.Handler,
			constants.PreReceiveHook:        prereceive.Handler,
			constants.PrepareCommitMsgHook:  preparecommitmsg.Handler,
			constants.UpdateHook:            update.Handler,
		},
	}
	fs.StringVar(&c.hook, "hook", "", "")
	return c
}

// Name returns command name
func (c *Command) Name() string {
	return c.fs.Name()
}

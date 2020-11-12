package handlers

import (
	"fisherman/clicontext"
	"fisherman/config"
	c "fisherman/constants"
)

type Handler interface {
	IsConfigured(config *config.HooksConfig) bool
	Handle(ctx *clicontext.CommandContext, args []string) error
}

var HandlerList = map[string]Handler{
	c.ApplyPatchMsgHook:     new(NotSupportedHandler),
	c.CommitMsgHook:         new(CommitMsgHandler),
	c.FsMonitorWatchmanHook: new(NotSupportedHandler),
	c.PostUpdateHook:        new(NotSupportedHandler),
	c.PreApplyPatchHook:     new(NotSupportedHandler),
	c.PreCommitHook:         new(PreCommitHandler),
	c.PrePushHook:           new(PrePushHandler),
	c.PreRebaseHook:         new(NotSupportedHandler),
	c.PreReceiveHook:        new(NotSupportedHandler),
	c.PrepareCommitMsgHook:  new(PrepareCommitMsgHandler),
	c.UpdateHook:            new(NotSupportedHandler),
}

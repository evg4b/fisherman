package handlers

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/handlers/shellhandlers"
	"fisherman/infrastructure/log"
)

// PrePushHandler is a handler for commit-msg hook
type PrePushHandler struct {
}

func (*PrePushHandler) IsConfigured(c *config.HooksConfig) bool {
	return c.PrePushHook.Variables != hooks.Variables{} || len(c.PrePushHook.Shell) > 0
}

// Handle is a handler for pre-push hook
func (*PrePushHandler) Handle(ctx *clicontext.CommandContext, args []string) error {
	config := ctx.Config.PrePushHook
	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		log.Debugf("Additional variables loading filed: %s", err)

		return err
	}

	config.Compile(ctx.Variables())

	return shellhandlers.ExecParallel(ctx, ctx.Shell, config.Shell)
}

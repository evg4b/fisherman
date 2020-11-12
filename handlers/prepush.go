package handlers

import (
	"fisherman/clicontext"
	"fisherman/handlers/shellhandlers"
	"fisherman/infrastructure/log"
)

// PrePushHandler is a handler for pre-push hook
func PrePushHandler(ctx *clicontext.CommandContext, args []string) error {
	config := ctx.Config.PrePushHook
	err := ctx.LoadAdditionalVariables(&config.Variables)
	if err != nil {
		log.Debugf("Additional variables loading filed: %s", err)

		return err
	}

	config.Compile(ctx.Variables())

	return shellhandlers.ExecParallel(ctx, ctx.Shell, config.Shell)
}

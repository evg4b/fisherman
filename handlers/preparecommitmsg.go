package handlers

import (
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/internal/clicontext"
	"fisherman/utils"
)

// PrePushHandler is a handler for commit-msg hook
type PrepareCommitMsgHandler struct {
}

func (*PrepareCommitMsgHandler) IsConfigured(c *config.HooksConfig) bool {
	return c.PrepareCommitMsgHook != hooks.PrepareCommitMsgHookConfig{}
}

// Handle is a handler for pre-push hook
func (*PrepareCommitMsgHandler) Handle(ctx *clicontext.CommandContext, args []string) error {
	config := &ctx.Config.PrepareCommitMsgHook
	if utils.IsNotEmpty(config.Message) {
		err := ctx.LoadAdditionalVariables(&config.Variables)
		if err != nil {
			return err
		}

		utils.FillTemplate(&config.Message, ctx.Variables())

		err = ctx.Files.Write(args[0], config.Message)
		utils.HandleCriticalError(err)
	}

	return nil
}

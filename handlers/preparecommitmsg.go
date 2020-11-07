package handlers

import (
	"fisherman/clicontext"
	"fisherman/utils"
)

// PrepareCommitMsgHandler is a execute function for prepare-commit-msg hook
func PrepareCommitMsgHandler(ctx *clicontext.CommandContext, args []string) error {
	config := &ctx.Config.PrepareCommitMsgHook
	if utils.IsNotEmpty(config.Message) {
		err := ctx.LoadAdditionalVariables(&config.Variables)
		if err != nil {
			return err
		}

		utils.FillTemplate(&config.Message, ctx.Variables)

		err = ctx.Files.Write(args[0], config.Message)
		utils.HandleCriticalError(err)
	}

	return nil
}

package handlers

import (
	"fisherman/commands/context"
	"fisherman/config/hooks"
	"fmt"
	"strings"
)

// CommitMsgHandler is a handler for commit-msg hook
func CommitMsgHandler(ctx *context.CommandContext, args []string) error {
	appConfig := ctx.GetHookConfiguration()
	if appConfig.CommitMsgHook == nil {
		return fmt.Errorf("")
	}

	commitMessage, err := ctx.FileAccessor.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("")
	}

	return validateMessage(commitMessage, appConfig.CommitMsgHook)
}

func validateMessage(message string, config *hooks.CommitMsgHookConfig) error {
	if config.NotEmpty {
		if len(strings.TrimSpace(message)) == 0 {
			return fmt.Errorf("demo")
		}
	}

	return nil
}

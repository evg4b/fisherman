package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PrePushHandler is a handler for pre-push hook
func PrePushHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

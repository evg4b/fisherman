package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// PrePushHandler is a handler for pre-push hook
func PrePushHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

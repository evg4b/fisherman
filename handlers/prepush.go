package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PrePushHandler is a handler for pre-push hook
func PrePushHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("pre-push")
	return nil
}

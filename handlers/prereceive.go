package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PreReceiveHandler is a handler for pre-receive hook
func PreReceiveHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("pre-receive")
	return nil
}

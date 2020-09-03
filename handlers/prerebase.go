package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PreRebaseHandler is a handler for pre-rebase hook
func PreRebaseHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("pre-rebase")
	return nil
}

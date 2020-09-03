package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PreApplyPatchHandler is a handler for pre-applypatch hook
func PreApplyPatchHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("pre-applypatch")
	return nil
}

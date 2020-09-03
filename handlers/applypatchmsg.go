package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// ApplyPatchMsgHandler is a handler for applypatch-msg hook
func ApplyPatchMsgHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("applypatch-msg")
	return nil
}

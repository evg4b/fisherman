package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// ApplyPatchMsgHandler is a handler for applypatch-msg hook
func ApplyPatchMsgHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

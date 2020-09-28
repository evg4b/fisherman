package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// ApplyPatchMsgHandler is a handler for applypatch-msg hook
func ApplyPatchMsgHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

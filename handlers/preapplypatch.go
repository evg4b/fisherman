package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PreApplyPatchHandler is a handler for pre-applypatch hook
func PreApplyPatchHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

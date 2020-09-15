package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// PreApplyPatchHandler is a handler for pre-applypatch hook
func PreApplyPatchHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

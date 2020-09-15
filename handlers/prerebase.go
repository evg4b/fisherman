package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// PreRebaseHandler is a handler for pre-rebase hook
func PreRebaseHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

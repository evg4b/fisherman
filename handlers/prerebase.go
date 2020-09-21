package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PreRebaseHandler is a handler for pre-rebase hook
func PreRebaseHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

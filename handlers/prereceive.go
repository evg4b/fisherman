package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PreReceiveHandler is a handler for pre-receive hook
func PreReceiveHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

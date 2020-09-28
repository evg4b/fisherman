package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PostUpdateHandler is a handler for post-update hook
func PostUpdateHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

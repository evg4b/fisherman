package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// PostUpdateHandler is a handler for post-update hook
func PostUpdateHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

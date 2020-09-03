package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PostUpdateHandler is a handler for post-update hook
func PostUpdateHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("post-update")
	return nil
}

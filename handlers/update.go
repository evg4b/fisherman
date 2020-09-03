package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// UpdateHandler is a handler for update hook
func UpdateHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("update")
	return nil
}

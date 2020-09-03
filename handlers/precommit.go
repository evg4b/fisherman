package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// PreCommitHandler is a handler for pre-commit hook
func PreCommitHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("pre-commit")
	return nil
}

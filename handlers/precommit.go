package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// PreCommitHandler is a handler for pre-commit hook
func PreCommitHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

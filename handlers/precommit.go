package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// PreCommitHandler is a handler for pre-commit hook
func PreCommitHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

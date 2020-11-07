package handlers

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fmt"
)

// PreRebaseHandler is a handler for pre-rebase hook
func PreRebaseHandler(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

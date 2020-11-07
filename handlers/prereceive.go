package handlers

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fmt"
)

// PreReceiveHandler is a handler for pre-receive hook
func PreReceiveHandler(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

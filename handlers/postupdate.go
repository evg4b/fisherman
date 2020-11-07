package handlers

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fmt"
)

// PostUpdateHandler is a handler for post-update hook
func PostUpdateHandler(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

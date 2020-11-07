package handlers

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fmt"
)

// ApplyPatchMsgHandler is a handler for applypatch-msg hook
func ApplyPatchMsgHandler(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

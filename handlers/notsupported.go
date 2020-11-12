package handlers

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/constants"
	"fmt"
)

// NotSupportedHandler is structure for a handler which is currently not supported
type NotSupportedHandler struct{}

// Handle is a handler for applypatch-msg hook
func (*NotSupportedHandler) Handle(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

// Handle is a handler for applypatch-msg hook
func (*NotSupportedHandler) IsConfigured(*config.HooksConfig) bool {
	return true
}

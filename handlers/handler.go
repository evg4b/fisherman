package handlers

import (
	"fisherman/commands/context"
)

// HookHandler is base handler interface
type HookHandler = func(ctx *context.CommandContext, args []string) error

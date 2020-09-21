package handlers

import "fisherman/commands"

// HookHandler is base handler interface
type HookHandler = func(ctx *commands.CommandContext, args []string) error

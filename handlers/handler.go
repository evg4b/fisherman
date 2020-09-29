package handlers

import "fisherman/commands"

// HookHandler is base handler interface
type HookHandler = func(ctx *commands.CommandContext, args []string) error

type HandlerFactory = func(ctx *commands.CommandContext, args []string) (HookHandler, error)

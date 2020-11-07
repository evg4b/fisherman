package handlers

import "fisherman/clicontext"

// HookHandler is base handler interface
type HookHandler = func(ctx *clicontext.CommandContext, args []string) error

type HandlerFactory = func(ctx *clicontext.CommandContext, args []string) (HookHandler, error)

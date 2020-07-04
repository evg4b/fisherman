package handler

import (
	"fisherman/commands/context"
)

// HookHandler is base handler interface
type HookHandler interface {
	Execute(ctx context.Context, args []string)
}

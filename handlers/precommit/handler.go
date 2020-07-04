package precommit

import (
	"fisherman/commands/context"
	"fmt"
)

type Handler struct {
}

// NewHandler is constructor for applypatch-msg hook handler
func NewHandler() *Handler {
	return &Handler{}
}

// Handler is a handler for pre-commit hook
func (h *Handler) Execute(ctx context.Context, args []string) {
	fmt.Print("pre-commit")
}

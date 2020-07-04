package applypatchmsg

import (
	"fisherman/commands/context"
	"fmt"
)

// Handler is a handler for applypatch-msg hook
type Handler struct {
}

// NewHandler is constructor for applypatch-msg hook handler
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Execute(ctx context.Context, args []string) {
	fmt.Print("applypatch-msg")
}

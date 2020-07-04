package applypatchmsg

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for applypatch-msg hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("applypatch-msg")
}

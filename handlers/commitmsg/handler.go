package commitmsg

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for commit-msg hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("commit-msg")
}

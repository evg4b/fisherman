package fsmonitorwatchman

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for fsmonitor-watchman hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("fsmonitor-watchman")
}

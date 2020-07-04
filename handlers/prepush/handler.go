package prepush

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for pre-push hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("pre-push")
}

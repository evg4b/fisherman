package prerebase

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for pre-rebase hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("pre-rebase")
}

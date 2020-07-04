package precommit

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for pre-commit hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("pre-commit")
}

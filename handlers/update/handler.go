package update

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for update hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("update")
}

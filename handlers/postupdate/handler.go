package postupdate

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for post-update hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("post-update")
}

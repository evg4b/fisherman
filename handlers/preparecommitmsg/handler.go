package preparecommitmsg

import (
	"fisherman/commands/context"
	appConfig "fisherman/config"
	"fmt"
)

// Handler is a handler for prepare-commit-msg
func Handler(ctx context.Context, config *appConfig.FishermanConfig) {
	fmt.Print("prepare-commit-msg")
}

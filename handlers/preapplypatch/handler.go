package preapplypatch

import (
	"fisherman/commands/context"
	"fisherman/config"
	"fmt"
)

// Handler is a handler for pre-applypatch hook
func Handler(ctx context.Context, config *config.FishermanConfig) {
	fmt.Print("pre-applypatch")
}

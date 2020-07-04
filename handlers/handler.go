package handler

import (
	"fisherman/commands/context"
	appConfig "fisherman/config"
)

// HookHandler is general type for hook handler
type HookHandler = func(ctx context.Context, config *appConfig.FishermanConfig)

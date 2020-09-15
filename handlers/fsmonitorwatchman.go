package handlers

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fmt"
)

// FsMonitorWatchmanHandler is a handler for fsmonitor-watchman hook
func FsMonitorWatchmanHandler(ctx *context.CommandContext, args []string) error {
	return fmt.Errorf("This hook is not supported in version %s.", constants.Version)
}

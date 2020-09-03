package handlers

import (
	"fisherman/commands/context"
	"fmt"
)

// FsMonitorWatchmanHandler is a handler for fsmonitor-watchman hook
func FsMonitorWatchmanHandler(ctx *context.CommandContext, args []string) error {
	fmt.Print("fsmonitor-watchman")
	return nil
}

package handlers

import (
	"fisherman/commands"
	"fisherman/constants"
	"fmt"
)

// FsMonitorWatchmanHandler is a handler for fsmonitor-watchman hook
func FsMonitorWatchmanHandler(ctx *commands.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

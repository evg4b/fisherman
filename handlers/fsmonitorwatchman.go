package handlers

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fmt"
)

// FsMonitorWatchmanHandler is a handler for fsmonitor-watchman hook
func FsMonitorWatchmanHandler(ctx *clicontext.CommandContext, args []string) error {
	return fmt.Errorf("this hook is not supported in version %s", constants.Version)
}

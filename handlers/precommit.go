package handlers

import (
	"fisherman/commands"
	"fisherman/handlers/common"
)

// PreCommitHandler is a handler for pre-commit hook
func PreCommitHandler(ctx *commands.CommandContext, args []string) error {
	ctx.Config.PreCommitHook.Compile(ctx.Variables)

	return common.ExecCommandsParallel(ctx.Shell, ctx.Config.PreCommitHook.Cmd)
}

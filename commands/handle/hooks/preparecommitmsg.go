package hooks

import (
	"fisherman/actions"
	"fisherman/config/hooks"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
)

func PrepareCommitMsg(
	ctxFactory internal.CtxFactory,
	configuration hooks.PrepareCommitMsgHookConfig,
	compile configcompiler.Compiler,
) *handling.HookHandler {
	compile(&configuration)

	return handling.NewHookHandler(
		ctxFactory,
		[]handling.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.PrepareMessage(ctx, configuration.Message)
			},
		},
		NoSyncValidators,
		NoAsyncValidators,
		NoAfterActions,
	)
}

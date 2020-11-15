package hooks

import (
	"fisherman/actions"
	c "fisherman/config/hooks"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	h "fisherman/internal/handling"
)

func PrepareCommitMsg(
	ctxFactory internal.CtxFactory,
	configuration c.PrepareCommitMsgHookConfig,
	compile configcompiler.Compiler,
) *h.HookHandler {
	compile(&configuration)

	return h.NewHookHandler(
		ctxFactory,
		[]h.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.PrepareMessage(ctx, configuration.Message)
			},
		},
		NoSyncValidators,
		NoAsyncValidators,
		NoAfterActions,
	)
}

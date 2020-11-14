package hooks

import (
	"fisherman/actions"
	c "fisherman/config/hooks"
	h "fisherman/internal/handling"
	v "fisherman/internal/validation"
)

func PrepareCommitMsg(factory ctxFactory, conf c.PrepareCommitMsgHookConfig, extr v.VariablesExtractor) *h.HookHandler {
	variables, err := extr.Variables(conf.Variables)
	if err != nil {
		panic(err)
	}

	conf.Compile(variables)

	return h.NewHookHandler(
		factory,
		[]h.Action{
			func(ctx v.SyncValidationContext) (bool, error) {
				return actions.PrepareMessage(ctx, conf.Message)
			},
		},
		NoSyncValidators,
		NoAsyncValidators,
		NoAfterActions,
	)
}

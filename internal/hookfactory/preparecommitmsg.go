package hookfactory

import (
	"fisherman/actions"
	hooks "fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/handling"
)

func (factory *Factory) PrepareCommitMsg(configuration hooks.PrepareCommitMsgHookConfig) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotRegistered
	}

	factory.compile(&configuration)

	return HandlerRegistration{
		Registered: true,
		Handler: handling.NewHookHandler(
			factory.ctxFactory,
			[]handling.Action{
				func(ctx internal.SyncContext) (bool, error) {
					return actions.PrepareMessage(ctx, configuration.Message)
				},
			},
			NoSyncValidators,
			NoAsyncValidators,
			NoAfterActions,
		),
	}
}

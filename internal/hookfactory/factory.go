package hookfactory

import (
	"fisherman/actions"
	hooks "fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/validators"
)

type Factory struct {
	ctxFactory internal.CtxFactory
	compile    configcompiler.Compiler
}

func NewFactory(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) *Factory {
	return &Factory{ctxFactory: ctxFactory, compile: compile}
}

func (factory *Factory) CommitMsg(configuration hooks.CommitMsgHookConfig) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotRegistered
	}

	factory.compile(&configuration)

	return HandlerRegistration{
		Registered: true,
		Handler: handling.NewHookHandler(
			factory.ctxFactory,
			NoBeforeActions,
			[]validation.SyncValidator{
				boolWrapper(validators.MessageNotEmpty, configuration.NotEmpty),
				stringWrapper(validators.MessageHasPrefix, configuration.MessagePrefix),
				stringWrapper(validators.MessageHasSuffix, configuration.MessageSuffix),
				stringWrapper(validators.MessageRegexp, configuration.MessageRegexp),
			},
			NoAsyncValidators,
			NoAfterActions,
		),
	}
}

func (factory *Factory) PreCommit(configuration hooks.PreCommitHookConfig) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotRegistered
	}

	factory.compile(&configuration)

	return HandlerRegistration{
		Registered: true,
		Handler: handling.NewHookHandler(
			factory.ctxFactory,
			NoBeforeActions,
			NoSyncValidators,
			scriptWrapper(configuration.Shell),
			[]handling.Action{
				func(ctx internal.SyncContext) (bool, error) {
					return actions.AddToIndex(ctx, configuration.AddFilesToIndex)
				},
			},
		),
	}
}

func (factory *Factory) PrePush(configuration hooks.PrePushHookConfig) HandlerRegistration {
	if configuration.IsEmpty() {
		return NotRegistered
	}

	factory.compile(&configuration)

	return HandlerRegistration{
		Registered: true,
		Handler: handling.NewHookHandler(
			factory.ctxFactory,
			NoBeforeActions,
			NoSyncValidators,
			scriptWrapper(configuration.Shell),
			NoAfterActions,
		),
	}
}

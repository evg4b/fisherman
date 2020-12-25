package hookfactory

import (
	"fisherman/actions"
	"fisherman/internal"
	"fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/validators"
)

func (factory *TFactory) commitMsg() *handling.HookHandler {
	configuration := factory.config.CommitMsgHook

	if configuration.IsEmpty() {
		return nil
	}

	factory.compile(&configuration)

	return &handling.HookHandler{
		ContextFactory: factory.ctxFactory,
		BeforeActions:  NoBeforeActions,
		SyncValidators: []validation.SyncValidator{
			boolWrapper(validators.MessageNotEmpty, configuration.NotEmpty),
			stringWrapper(validators.MessageHasPrefix, configuration.MessagePrefix),
			stringWrapper(validators.MessageHasSuffix, configuration.MessageSuffix),
			stringWrapper(validators.MessageRegexp, configuration.MessageRegexp),
		},
		AsyncValidators: NoAsyncValidators,
		AfterActions:    NoAfterActions,
	}
}

func (factory *TFactory) preCommit() *handling.HookHandler {
	configuration := factory.config.PreCommitHook
	if configuration.IsEmpty() {
		return nil
	}

	factory.compile(&configuration)

	return &handling.HookHandler{
		ContextFactory:  factory.ctxFactory,
		BeforeActions:   NoBeforeActions,
		SyncValidators:  NoSyncValidators,
		AsyncValidators: scriptWrapper(configuration.Shell),
		AfterActions: []handling.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.AddToIndex(ctx, configuration.AddFilesToIndex)
			},
			func(ctx internal.SyncContext) (bool, error) {
				return actions.SuppresCommitFiles(ctx, configuration.SuppressCommitFiles)
			},
		},
	}
}

func (factory *TFactory) prePush() *handling.HookHandler {
	configuration := factory.config.PrePushHook
	if configuration.IsEmpty() {
		return nil
	}

	factory.compile(&configuration)

	return &handling.HookHandler{
		ContextFactory:  factory.ctxFactory,
		BeforeActions:   NoBeforeActions,
		SyncValidators:  NoSyncValidators,
		AsyncValidators: scriptWrapper(configuration.Shell),
		AfterActions:    NoAfterActions,
	}
}

func (factory *TFactory) prepareCommitMsg() *handling.HookHandler {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration.IsEmpty() {
		return nil
	}

	factory.compile(&configuration)

	return &handling.HookHandler{
		ContextFactory: factory.ctxFactory,
		BeforeActions: []handling.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.PrepareMessage(ctx, configuration.Message)
			},
		},
		SyncValidators:  NoSyncValidators,
		AsyncValidators: NoAsyncValidators,
		AfterActions:    NoAfterActions,
	}
}

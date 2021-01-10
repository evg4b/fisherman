package hookfactory

import (
	"fisherman/actions"
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
	"fisherman/internal/validation"
	"fisherman/validators"
)

func (factory *TFactory) commitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, nil
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

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
	}, nil
}

func (factory *TFactory) preCommit() (*handling.HookHandler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, nil
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		ContextFactory:  factory.ctxFactory,
		BeforeActions:   NoBeforeActions,
		SyncValidators:  NoSyncValidators,
		AsyncValidators: scriptWrapper(configuration.Shell, expression.NewExpressionEngine(variables)),
		AfterActions: []handling.Action{
			func(ctx internal.SyncContext) (bool, error) {
				return actions.AddToIndex(ctx, configuration.AddFilesToIndex)
			},
			func(ctx internal.SyncContext) (bool, error) {
				return actions.SuppresCommitFiles(ctx, configuration.SuppressCommitFiles)
			},
		},
	}, nil
}

func (factory *TFactory) prePush() (*handling.HookHandler, error) {
	configuration := factory.config.PrePushHook
	if configuration == nil {
		return nil, nil
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		ContextFactory:  factory.ctxFactory,
		BeforeActions:   NoBeforeActions,
		SyncValidators:  NoSyncValidators,
		AsyncValidators: scriptWrapper(configuration.Shell, expression.NewExpressionEngine(variables)),
		AfterActions:    NoAfterActions,
	}, nil
}

func (factory *TFactory) prepareCommitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration == nil {
		return nil, nil
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

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
	}, nil
}

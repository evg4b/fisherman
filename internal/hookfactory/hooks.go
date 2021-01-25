package hookfactory

import (
	"fisherman/actions"
	"fisherman/internal"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
)

// TODO: move to configuration
const workersCount = 5

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
		BeforeActions:   NoBeforeActions,
		Rules:           getBaseRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		AsyncValidators: NoAsyncValidators,
		AfterActions:    NoAfterActions,
		WorkersCount:    workersCount,
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
		BeforeActions:   NoBeforeActions,
		AsyncValidators: scriptWrapper(configuration.Shell, expression.NewExpressionEngine(variables)),
		AfterActions: []handling.Action{
			func(ctx internal.ExecutionContext) (bool, error) {
				return actions.AddToIndex(ctx, configuration.AddFilesToIndex)
			},
			func(ctx internal.ExecutionContext) (bool, error) {
				return actions.SuppresCommitFiles(ctx, configuration.SuppressCommitFiles)
			},
		},
		WorkersCount: workersCount,
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
		BeforeActions:   NoBeforeActions,
		AsyncValidators: scriptWrapper(configuration.Shell, expression.NewExpressionEngine(variables)),
		AfterActions:    NoAfterActions,
		WorkersCount:    workersCount,
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
		BeforeActions: []handling.Action{
			func(ctx internal.ExecutionContext) (bool, error) {
				return actions.PrepareMessage(ctx, configuration.Message)
			},
		},
		AsyncValidators: NoAsyncValidators,
		AfterActions:    NoAfterActions,
		WorkersCount:    workersCount,
	}, nil
}

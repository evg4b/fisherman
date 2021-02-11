package hookfactory

import (
	"errors"
	"fisherman/internal/handling"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

func (factory *TFactory) commitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScripts(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *TFactory) preCommit() (*handling.HookHandler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScripts(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *TFactory) prePush() (*handling.HookHandler, error) {
	configuration := factory.config.PrePushHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScripts(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *TFactory) prepareCommitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	variables, err := factory.prepareConfig(configuration)
	if err != nil || variables == nil {
		return nil, err
	}

	return &handling.HookHandler{
		WorkersCount: workersCount,
	}, nil
}

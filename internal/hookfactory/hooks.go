package hookfactory

import (
	"errors"
	"fisherman/internal/handling"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

func (factory *GitHookFactory) commitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := factory.prepareConfig(configuration)
	if err != nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScriptRules(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *GitHookFactory) preCommit() (*handling.HookHandler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := factory.prepareConfig(configuration)
	if err != nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScriptRules(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *GitHookFactory) prePush() (*handling.HookHandler, error) {
	configuration := factory.config.PrePushHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := factory.prepareConfig(configuration)
	if err != nil {
		return nil, err
	}

	return &handling.HookHandler{
		Rules:           getPreScriptRules(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *GitHookFactory) prepareCommitMsg() (*handling.HookHandler, error) {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := factory.prepareConfig(configuration)
	if err != nil {
		return nil, err
	}

	return &handling.HookHandler{
		WorkersCount: workersCount,
	}, nil
}

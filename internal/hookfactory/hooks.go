package hookfactory

import (
	"errors"
	"fisherman/configuration"
	"fisherman/internal/handling"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

func (factory *GitHookFactory) commitMsg() (handling.Handler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
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

func (factory *GitHookFactory) preCommit() (handling.Handler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
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

func (factory *GitHookFactory) prePush() (handling.Handler, error) {
	configuration := factory.config.PrePushHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) prepareCommitMsg() (handling.Handler, error) {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &handling.HookHandler{
		WorkersCount: workersCount,
	}, nil
}

func (factory *GitHookFactory) applyPatchMsg() (handling.Handler, error) {
	configuration := factory.config.ApplyPatchMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) fsMonitorWatchman() (handling.Handler, error) {
	configuration := factory.config.FsMonitorWatchmanHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) postUpdate() (handling.Handler, error) {
	configuration := factory.config.PostUpdateHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) preApplyPatch() (handling.Handler, error) {
	configuration := factory.config.PreApplyPatchHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) preRebase() (handling.Handler, error) {
	configuration := factory.config.PreRebaseHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) preReceive() (handling.Handler, error) {
	configuration := factory.config.PreReceiveHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) update() (handling.Handler, error) {
	configuration := factory.config.UpdateHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return factory.configureCommon(&configuration.CommonConfig), nil
}

func (factory *GitHookFactory) configureCommon(config *configuration.CommonConfig) *handling.HookHandler {
	return &handling.HookHandler{
		Rules:           getPreScriptRules(config.Rules),
		Scripts:         getScriptRules(config.Rules),
		PostScriptRules: getPostScriptRules(config.Rules),
		WorkersCount:    workersCount,
	}
}

package handling

import (
	"errors"
	"fisherman/configuration"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

func (factory *GitHookFactory) commitMsg() (Handler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &HookHandler{
		Rules:           getPreScriptRules(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *GitHookFactory) preCommit() (Handler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &HookHandler{
		Rules:           getPreScriptRules(configuration.Rules),
		Scripts:         getScriptRules(configuration.Rules),
		PostScriptRules: getPostScriptRules(configuration.Rules),
		WorkersCount:    workersCount,
	}, nil
}

func (factory *GitHookFactory) prePush() (Handler, error) {
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

func (factory *GitHookFactory) prepareCommitMsg() (Handler, error) {
	configuration := factory.config.PrepareCommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	err := configuration.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	return &HookHandler{
		WorkersCount: workersCount,
	}, nil
}

func (factory *GitHookFactory) applyPatchMsg() (Handler, error) {
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

func (factory *GitHookFactory) fsMonitorWatchman() (Handler, error) {
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

func (factory *GitHookFactory) postUpdate() (Handler, error) {
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

func (factory *GitHookFactory) preApplyPatch() (Handler, error) {
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

func (factory *GitHookFactory) preRebase() (Handler, error) {
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

func (factory *GitHookFactory) preReceive() (Handler, error) {
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

func (factory *GitHookFactory) update() (Handler, error) {
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

func (factory *GitHookFactory) configureCommon(config *configuration.CommonConfig) *HookHandler {
	return &HookHandler{
		Rules:           getPreScriptRules(config.Rules),
		Scripts:         getScriptRules(config.Rules),
		PostScriptRules: getPostScriptRules(config.Rules),
		WorkersCount:    workersCount,
	}
}

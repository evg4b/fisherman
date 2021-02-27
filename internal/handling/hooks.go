package handling

import (
	"errors"
	"fisherman/configuration"
	"fisherman/internal/rules"
	"fisherman/utils"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

func (factory *GitHookFactory) commitMsg() (Handler, error) {
	configuration := factory.config.CommitMsgHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("commit-msg hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) preCommit() (Handler, error) {
	configuration := factory.config.PreCommitHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.AddToIndexType,
		rules.CommitMessageType,
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("pre-commit hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) prePush() (Handler, error) {
	configuration := factory.config.PrePushHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("pre-push hook: %v", err)
	}

	return handler, nil
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

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("apply-patch-msg hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) fsMonitorWatchman() (Handler, error) {
	configuration := factory.config.FsMonitorWatchmanHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("fs-monitor-watchman hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) postUpdate() (Handler, error) {
	configuration := factory.config.PostUpdateHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("post-update hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) preApplyPatch() (Handler, error) {
	configuration := factory.config.PreApplyPatchHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.AddToIndexType,
		rules.CommitMessageType,
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("pre-apply-patch hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) preRebase() (Handler, error) {
	configuration := factory.config.PreRebaseHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("pre-rebase hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) preReceive() (Handler, error) {
	configuration := factory.config.PreReceiveHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("pre-receive hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) update() (Handler, error) {
	configuration := factory.config.UpdateHook
	if configuration == nil {
		return nil, ErrNotPresented
	}

	handler, err := factory.configureCommon(&configuration.HookConfig, []string{
		rules.ShellScriptType,
	})

	if err != nil {
		return nil, fmt.Errorf("update hook: %v", err)
	}

	return handler, nil
}

func (factory *GitHookFactory) configureCommon(
	config *configuration.HookConfig,
	allowed []string,
) (*HookHandler, error) {
	err := config.Compile(factory.engine, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var multiError *multierror.Error
	for _, rule := range config.Rules {
		if !utils.Contains(allowed, rule.GetType()) {
			multiError = multierror.Append(multiError, fmt.Errorf("rule %s is not allowed", rule.GetType()))
		}
	}

	return &HookHandler{
		Rules:           getPreScriptRules(config.Rules),
		Scripts:         getScriptRules(config.Rules),
		PostScriptRules: getPostScriptRules(config.Rules),
		WorkersCount:    workersCount,
	}, multiError.ErrorOrNil()
}

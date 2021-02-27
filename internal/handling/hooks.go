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
	return factory.configureCommon(factory.config.CommitMsgHook, []string{
		rules.AddToIndexType,
		rules.CommitMessageType,
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) preCommit() (Handler, error) {
	return factory.configureCommon(factory.config.PreCommitHook, []string{
		rules.AddToIndexType,
		rules.CommitMessageType,
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) prePush() (Handler, error) {
	return factory.configureCommon(factory.config.PrePushHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) prepareCommitMsg() (Handler, error) {
	return factory.configureCommon(factory.config.PrepareCommitMsgHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) applyPatchMsg() (Handler, error) {
	return factory.configureCommon(factory.config.ApplyPatchMsgHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) fsMonitorWatchman() (Handler, error) {
	return factory.configureCommon(factory.config.FsMonitorWatchmanHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) postUpdate() (Handler, error) {
	return factory.configureCommon(factory.config.PostUpdateHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) preApplyPatch() (Handler, error) {
	return factory.configureCommon(factory.config.PreApplyPatchHook, []string{
		rules.AddToIndexType,
		rules.CommitMessageType,
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) preRebase() (Handler, error) {
	return factory.configureCommon(factory.config.PreRebaseHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) preReceive() (Handler, error) {
	return factory.configureCommon(factory.config.PreReceiveHook, []string{
		rules.ShellScriptType,
	})
}

func (factory *GitHookFactory) update() (Handler, error) {
	return factory.configureCommon(
		factory.config.UpdateHook,
		[]string{rules.ShellScriptType})
}

func (factory *GitHookFactory) configureCommon(
	config *configuration.HookConfig,
	allowed []string,
) (*HookHandler, error) {
	if config == nil {
		return nil, ErrNotPresented
	}

	// TODO: Provide hook name
	name := "demo"

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

	err = multiError.ErrorOrNil()
	if err != nil {
		return nil, fmt.Errorf("%s hook: %v", name, err)
	}

	return &HookHandler{
		Rules:           getPreScriptRules(config.Rules),
		Scripts:         getScriptRules(config.Rules),
		PostScriptRules: getPostScriptRules(config.Rules),
		WorkersCount:    workersCount,
	}, multiError.ErrorOrNil()
}

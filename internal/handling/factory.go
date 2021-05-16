package handling

import (
	"errors"
	cnfg "fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/expression"
	"fisherman/internal/rules"
	"fisherman/utils"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration
const workersCount = 5

type CompilableConfig interface {
	Compile(engine expression.Engine, global Variables) (Variables, error)
}

type Factory interface {
	GetHook(name string) (Handler, error)
}

type Variables = map[string]interface{}
type hookBuilder = func() (Handler, error)
type builders = map[string]hookBuilder

type GitHookFactory struct {
	engine        expression.Engine
	hooksBuilders builders
}

func NewFactory(engine expression.Engine, config cnfg.HooksConfig) *GitHookFactory {
	factory := GitHookFactory{
		engine: engine,
	}

	factory.hooksBuilders = builders{
		constants.ApplyPatchMsgHook: factory.configure(
			constants.ApplyPatchMsgHook,
			config.ApplyPatchMsgHook,
			[]string{rules.ShellScriptType},
		),
		constants.CommitMsgHook: factory.configure(
			constants.CommitMsgHook,
			config.CommitMsgHook,
			[]string{
				rules.ShellScriptType,
				rules.CommitMessageType,
			},
		),
		constants.FsMonitorWatchmanHook: factory.configure(
			constants.FsMonitorWatchmanHook,
			config.FsMonitorWatchmanHook,
			[]string{rules.ShellScriptType},
		),
		constants.PostUpdateHook: factory.configure(
			constants.PostUpdateHook,
			config.PostUpdateHook,
			[]string{rules.ShellScriptType},
		),
		constants.PreApplyPatchHook: factory.configure(
			constants.PreApplyPatchHook,
			config.PreApplyPatchHook,
			[]string{rules.ShellScriptType},
		),
		constants.PreCommitHook: factory.configure(
			constants.PreCommitHook,
			config.PreCommitHook,
			[]string{
				rules.ShellScriptType,
				rules.AddToIndexType,
				rules.SuppressCommitFilesType,
			},
		),
		constants.PrePushHook: factory.configure(
			constants.PrePushHook,
			config.PrePushHook,
			[]string{rules.ShellScriptType},
		),
		constants.PreRebaseHook: factory.configure(
			constants.PreRebaseHook,
			config.PreRebaseHook,
			[]string{rules.ShellScriptType},
		),
		constants.PreReceiveHook: factory.configure(
			constants.PreReceiveHook,
			config.PreReceiveHook,
			[]string{rules.ShellScriptType},
		),
		constants.PrepareCommitMsgHook: factory.configure(
			constants.PrepareCommitMsgHook,
			config.PrepareCommitMsgHook,
			[]string{rules.PrepareMessageType},
		),
		constants.UpdateHook: factory.configure(
			constants.UpdateHook,
			config.UpdateHook,
			[]string{rules.ShellScriptType},
		),
	}

	return &factory
}

func (factory *GitHookFactory) GetHook(name string) (Handler, error) {
	if builder, ok := factory.hooksBuilders[name]; ok {
		return builder()
	}

	return nil, errors.New("unknown hook")
}

func (factory *GitHookFactory) configure(name string, config *cnfg.HookConfig, allowed []string) hookBuilder {
	return func() (Handler, error) {
		if config == nil {
			return nil, ErrNotPresented
		}

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
			Engine:          factory.engine,
			Rules:           getPreScriptRules(config.Rules),
			Scripts:         getScriptRules(config.Rules),
			PostScriptRules: getPostScriptRules(config.Rules),
			WorkersCount:    workersCount,
		}, nil
	}
}

package handling

import (
	"errors"
	"fisherman/configuration"
	c "fisherman/internal/constants"
	"fisherman/internal/expression"
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
	GetHook(name string, global map[string]interface{}) (Handler, error)
}

type Variables = map[string]interface{}
type hookBuilder = func(globalVars Variables) (Handler, error)
type builders = map[string]hookBuilder

type HookHandlerFactory struct {
	engine        expression.Engine
	hooksBuilders builders
}

func NewHookHandlerFactory(engine expression.Engine, config configuration.HooksConfig) *HookHandlerFactory {
	f := HookHandlerFactory{
		engine: engine,
	}

	f.hooksBuilders = builders{
		c.ApplyPatchMsgHook:     f.configure(c.ApplyPatchMsgHook, config.ApplyPatchMsgHook),
		c.CommitMsgHook:         f.configure(c.CommitMsgHook, config.CommitMsgHook),
		c.FsMonitorWatchmanHook: f.configure(c.FsMonitorWatchmanHook, config.FsMonitorWatchmanHook),
		c.PostUpdateHook:        f.configure(c.PostUpdateHook, config.PostUpdateHook),
		c.PreApplyPatchHook:     f.configure(c.PreApplyPatchHook, config.PreApplyPatchHook),
		c.PreCommitHook:         f.configure(c.PreCommitHook, config.PreCommitHook),
		c.PrePushHook:           f.configure(c.PrePushHook, config.PrePushHook),
		c.PreRebaseHook:         f.configure(c.PreRebaseHook, config.PreRebaseHook),
		c.PreReceiveHook:        f.configure(c.PreReceiveHook, config.PreReceiveHook),
		c.PrepareCommitMsgHook:  f.configure(c.PrepareCommitMsgHook, config.PrepareCommitMsgHook),
		c.UpdateHook:            f.configure(c.UpdateHook, config.UpdateHook),
	}

	return &f
}

func (f *HookHandlerFactory) GetHook(name string, globalVars Variables) (Handler, error) {
	if builder, ok := f.hooksBuilders[name]; ok {
		return builder(globalVars)
	}

	return nil, errors.New("unknown hook")
}

func (f *HookHandlerFactory) configure(name string, config *configuration.HookConfig) hookBuilder {
	return func(globalVars Variables) (Handler, error) {
		if config == nil {
			return nil, ErrNotPresented
		}

		compiledVars, err := config.Compile(globalVars)
		if err != nil {
			return nil, err
		}

		var multiError *multierror.Error
		for _, rule := range config.Rules {
			if !utils.Contains(allowedHooks[name], rule.GetType()) {
				multiError = multierror.Append(multiError, fmt.Errorf("rule %s is not allowed", rule.GetType()))
			}
		}

		err = multiError.ErrorOrNil()
		if err != nil {
			return nil, fmt.Errorf("%s hook: %v", name, err)
		}

		return &HookHandler{
			Engine:          f.engine,
			Rules:           getPreScriptRules(config.Rules),
			Scripts:         getScriptRules(config.Rules),
			PostScriptRules: getPostScriptRules(config.Rules),
			WorkersCount:    workersCount,
			GlobalVariables: compiledVars,
		}, nil
	}
}

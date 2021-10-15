package handling

import (
	"fisherman/internal/configuration"
	c "fisherman/internal/constants"
	"fisherman/internal/expression"
	"fisherman/internal/utils"

	"github.com/go-errors/errors"

	"github.com/hashicorp/go-multierror"
)

var ErrNotPresented = errors.New("configuration for hook is not presented")

// TODO: move to configuration.
const workersCount = 5

type CompilableConfig interface {
	Compile(engine expression.Engine, global Variables) (Variables, error)
}

type Factory interface {
	GetHook(name string, global map[string]interface{}) (Handler, error)
}

type (
	Variables   = map[string]interface{}
	hookBuilder = func(globalVars Variables) (Handler, error)
	builders    = map[string]hookBuilder
)

type HookHandlerFactory struct {
	engine        expression.Engine
	hooksBuilders builders
}

func NewHookHandlerFactory(engine expression.Engine, config configuration.HooksConfig) *HookHandlerFactory {
	return &HookHandlerFactory{
		engine: engine,
		hooksBuilders: builders{
			c.ApplyPatchMsgHook:     configure(engine, c.ApplyPatchMsgHook, config.ApplyPatchMsgHook),
			c.CommitMsgHook:         configure(engine, c.CommitMsgHook, config.CommitMsgHook),
			c.FsMonitorWatchmanHook: configure(engine, c.FsMonitorWatchmanHook, config.FsMonitorWatchmanHook),
			c.PostUpdateHook:        configure(engine, c.PostUpdateHook, config.PostUpdateHook),
			c.PreApplyPatchHook:     configure(engine, c.PreApplyPatchHook, config.PreApplyPatchHook),
			c.PreCommitHook:         configure(engine, c.PreCommitHook, config.PreCommitHook),
			c.PrePushHook:           configure(engine, c.PrePushHook, config.PrePushHook),
			c.PreRebaseHook:         configure(engine, c.PreRebaseHook, config.PreRebaseHook),
			c.PreReceiveHook:        configure(engine, c.PreReceiveHook, config.PreReceiveHook),
			c.PrepareCommitMsgHook:  configure(engine, c.PrepareCommitMsgHook, config.PrepareCommitMsgHook),
			c.UpdateHook:            configure(engine, c.UpdateHook, config.UpdateHook),
		},
	}
}

func (f *HookHandlerFactory) GetHook(name string, globalVars Variables) (Handler, error) {
	if builder, ok := f.hooksBuilders[name]; ok {
		return builder(globalVars)
	}

	return nil, errors.New("unknown hook")
}

func configure(engine expression.Engine, name string, config *configuration.HookConfig) hookBuilder {
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
				multiError = multierror.Append(multiError, errors.Errorf("rule %s is not allowed", rule.GetType()))
			}
		}

		err = multiError.ErrorOrNil()
		if err != nil {
			return nil, errors.Errorf("%s hook: %v", name, err)
		}

		return &HookHandler{
			Engine:          engine,
			Rules:           getPreScriptRules(config.Rules),
			Scripts:         getScriptRules(config.Rules),
			PostScriptRules: getPostScriptRules(config.Rules),
			WorkersCount:    workersCount,
			GlobalVariables: compiledVars,
		}, nil
	}
}

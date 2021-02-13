package hookfactory

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/expression"
	"fisherman/internal/handling"
)

type CompilableConfig interface {
	Compile(engine expression.Engine, global map[string]interface{})
	GetVariables() map[string]interface{}
}

type builders = map[string]func() (*handling.HookHandler, error)

type Factory interface {
	GetHook(name string) (handling.Handler, error)
}

type GitHookFactory struct {
	engine        expression.Engine
	config        configuration.HooksConfig
	hooksBuilders builders
}

func NewFactory(engine expression.Engine, config configuration.HooksConfig) *GitHookFactory {
	factory := GitHookFactory{
		engine: engine,
		config: config,
	}

	factory.hooksBuilders = builders{
		constants.CommitMsgHook:        factory.commitMsg,
		constants.PreCommitHook:        factory.preCommit,
		constants.PrePushHook:          factory.prePush,
		constants.PrepareCommitMsgHook: factory.prepareCommitMsg,
	}

	return &factory
}

func (factory *GitHookFactory) GetHook(name string) (handling.Handler, error) {
	if builder, ok := factory.hooksBuilders[name]; ok {
		return builder()
	}

	return nil, errors.New("unknown hook")
}

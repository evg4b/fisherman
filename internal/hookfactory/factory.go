package hookfactory

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
)

type builders = map[string]func() (*handling.HookHandler, error)

type Factory interface {
	GetHook(name string) (handling.Handler, error)
}

type TFactory struct {
	extractor     configcompiler.Extractor
	config        configuration.HooksConfig
	hooksBuilders builders
}

func NewFactory(extractor configcompiler.Extractor, config configuration.HooksConfig) *TFactory {
	factory := TFactory{
		extractor: extractor,
		config:    config,
	}

	factory.hooksBuilders = builders{
		constants.CommitMsgHook:        factory.commitMsg,
		constants.PreCommitHook:        factory.preCommit,
		constants.PrePushHook:          factory.prePush,
		constants.PrepareCommitMsgHook: factory.prepareCommitMsg,
	}

	return &factory
}

func (factory *TFactory) GetHook(name string) (handling.Handler, error) {
	if builder, ok := factory.hooksBuilders[name]; ok {
		return builder()
	}

	return nil, errors.New("unknown hook")
}

func (factory *TFactory) prepareConfig(configuration configcompiler.CompilableConfig) (map[string]interface{}, error) {
	if configuration == nil {
		return nil, nil
	}

	variables, err := factory.extractor.Variables(configuration.GetVariablesConfig())
	if err != nil {
		return nil, err
	}

	configuration.Compile(variables)

	return variables, nil
}

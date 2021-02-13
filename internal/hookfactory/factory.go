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

type GitHookFactory struct {
	extractor     configcompiler.Extractor
	config        configuration.HooksConfig
	hooksBuilders builders
}

func NewFactory(extractor configcompiler.Extractor, config configuration.HooksConfig) *GitHookFactory {
	factory := GitHookFactory{
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

func (factory *GitHookFactory) GetHook(name string) (handling.Handler, error) {
	if builder, ok := factory.hooksBuilders[name]; ok {
		return builder()
	}

	return nil, errors.New("unknown hook")
}

func (factory *GitHookFactory) prepareConfig(configuration configcompiler.CompilableConfig) error {
	variables, err := factory.extractor.Variables(configuration.GetVariablesConfig())
	if err != nil {
		return err
	}

	configuration.Compile(variables)

	return nil
}

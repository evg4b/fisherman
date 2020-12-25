package hookfactory

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
)

type builders = map[string]func() *handling.HookHandler

type Factory interface {
	GetHook(name string) (handling.Handler, error)
}

type TFactory struct {
	ctxFactory    internal.CtxFactory
	compile       configcompiler.Compiler
	config        configuration.HooksConfig
	hooksBuilders builders
}

func NewFactory(
	ctxFactory internal.CtxFactory,
	compile configcompiler.Compiler,
	config configuration.HooksConfig,
) *TFactory {
	factory := TFactory{
		ctxFactory: ctxFactory,
		compile:    compile,
		config:     config,
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
		return builder(), nil
	}

	return nil, errors.New("unknown hook")
}

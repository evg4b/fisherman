package handling

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/expression"
)

type Variables = map[string]interface{}
type CompilableConfig interface {
	Compile(engine expression.Engine, global Variables) (Variables, error)
}

type builders = map[string]func() (Handler, error)

type Factory interface {
	GetHook(name string) (Handler, error)
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
		constants.ApplyPatchMsgHook:     factory.applyPatchMsg,
		constants.CommitMsgHook:         factory.commitMsg,
		constants.FsMonitorWatchmanHook: factory.fsMonitorWatchman,
		constants.PostUpdateHook:        factory.postUpdate,
		constants.PreApplyPatchHook:     factory.preApplyPatch,
		constants.PreCommitHook:         factory.preCommit,
		constants.PrePushHook:           factory.prePush,
		constants.PreRebaseHook:         factory.preRebase,
		constants.PreReceiveHook:        factory.preReceive,
		constants.PrepareCommitMsgHook:  factory.prepareCommitMsg,
		constants.UpdateHook:            factory.update,
	}

	return &factory
}

func (factory *GitHookFactory) GetHook(name string) (Handler, error) {
	if builder, ok := factory.hooksBuilders[name]; ok {
		return builder()
	}

	return nil, errors.New("unknown hook")
}

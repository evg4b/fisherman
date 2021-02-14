package hookfactory_test

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/hookfactory"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetHook(t *testing.T) {
	factory := hookfactory.NewFactory(
		mocks.NewEngineMock(t).EvalMock.Return(false, nil),
		configuration.HooksConfig{
			ApplyPatchMsgHook:     &configuration.ApplyPatchMsgHookConfig{},
			FsMonitorWatchmanHook: &configuration.FsMonitorWatchmanHookConfig{},
			PostUpdateHook:        &configuration.PostUpdateHookConfig{},
			PreApplyPatchHook:     &configuration.PreApplyPatchHookConfig{},
			PreCommitHook:         &configuration.PreCommitHookConfig{},
			PrePushHook:           &configuration.PrePushHookConfig{},
			PreRebaseHook:         &configuration.PreRebaseHookConfig{},
			PreReceiveHook:        &configuration.PreReceiveHookConfig{},
			UpdateHook:            &configuration.UpdateHookConfig{},
			CommitMsgHook:         &configuration.CommitMsgHookConfig{},
			PrepareCommitMsgHook:  &configuration.PrepareCommitMsgHookConfig{Message: "test"},
		},
	)

	for _, tt := range constants.HooksNames {
		t.Run(tt, func(t *testing.T) {
			hook, err := factory.GetHook(tt)

			assert.NotNil(t, hook)
			assert.NoError(t, err)
		})
	}
}

func TestFactory_GetHook_ReturnInternalError(t *testing.T) {
	variablesSection := configuration.VariablesSection{ExtractVariables: []string{"stub"}}
	commonConfig := configuration.CommonConfig{VariablesSection: variablesSection}

	factory := hookfactory.NewFactory(
		mocks.NewEngineMock(t).EvalMapMock.Return(nil, errors.New("test error")),
		configuration.HooksConfig{
			ApplyPatchMsgHook:     &configuration.ApplyPatchMsgHookConfig{CommonConfig: commonConfig},
			FsMonitorWatchmanHook: &configuration.FsMonitorWatchmanHookConfig{CommonConfig: commonConfig},
			PostUpdateHook:        &configuration.PostUpdateHookConfig{CommonConfig: commonConfig},
			PreApplyPatchHook:     &configuration.PreApplyPatchHookConfig{CommonConfig: commonConfig},
			PreCommitHook:         &configuration.PreCommitHookConfig{CommonConfig: commonConfig},
			PrePushHook:           &configuration.PrePushHookConfig{CommonConfig: commonConfig},
			PreRebaseHook:         &configuration.PreRebaseHookConfig{CommonConfig: commonConfig},
			PreReceiveHook:        &configuration.PreReceiveHookConfig{CommonConfig: commonConfig},
			UpdateHook:            &configuration.UpdateHookConfig{CommonConfig: commonConfig},
			CommitMsgHook:         &configuration.CommitMsgHookConfig{CommonConfig: commonConfig},
			PrepareCommitMsgHook:  &configuration.PrepareCommitMsgHookConfig{VariablesSection: variablesSection},
		},
	)

	for _, tt := range constants.HooksNames {
		t.Run(tt, func(t *testing.T) {
			hook, err := factory.GetHook(tt)

			assert.Nil(t, hook)
			assert.Error(t, err, "test error")
		})
	}
}

func TestFactory_GetHook_NotConfigured(t *testing.T) {
	factory := hookfactory.NewFactory(
		mocks.NewEngineMock(t),
		configuration.HooksConfig{},
	)

	for _, tt := range constants.HooksNames {
		t.Run(tt, func(t *testing.T) {
			hook, err := factory.GetHook(tt)

			assert.Nil(t, hook)
			assert.Equal(t, hookfactory.ErrNotPresented, err)
		})
	}
}

func TestFactory_GetHook_UnknownHook(t *testing.T) {
	factory := hookfactory.NewFactory(
		mocks.NewEngineMock(t),
		configuration.HooksConfig{},
	)

	hook, err := factory.GetHook("unknown-hook")

	assert.Nil(t, hook)
	assert.EqualError(t, err, "unknown hook")
}

package handling_test

import (
	"errors"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/handling"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetHook(t *testing.T) {
	factory := handling.NewFactory(
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
	commonConfig := configuration.HookConfig{VariablesSection: variablesSection}

	factory := handling.NewFactory(
		mocks.NewEngineMock(t).EvalMapMock.Return(nil, errors.New("test error")),
		configuration.HooksConfig{
			ApplyPatchMsgHook:     &configuration.ApplyPatchMsgHookConfig{HookConfig: commonConfig},
			FsMonitorWatchmanHook: &configuration.FsMonitorWatchmanHookConfig{HookConfig: commonConfig},
			PostUpdateHook:        &configuration.PostUpdateHookConfig{HookConfig: commonConfig},
			PreApplyPatchHook:     &configuration.PreApplyPatchHookConfig{HookConfig: commonConfig},
			PreCommitHook:         &configuration.PreCommitHookConfig{HookConfig: commonConfig},
			PrePushHook:           &configuration.PrePushHookConfig{HookConfig: commonConfig},
			PreRebaseHook:         &configuration.PreRebaseHookConfig{HookConfig: commonConfig},
			PreReceiveHook:        &configuration.PreReceiveHookConfig{HookConfig: commonConfig},
			UpdateHook:            &configuration.UpdateHookConfig{HookConfig: commonConfig},
			CommitMsgHook:         &configuration.CommitMsgHookConfig{HookConfig: commonConfig},
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
	factory := handling.NewFactory(
		mocks.NewEngineMock(t),
		configuration.HooksConfig{},
	)

	for _, tt := range constants.HooksNames {
		t.Run(tt, func(t *testing.T) {
			hook, err := factory.GetHook(tt)

			assert.Nil(t, hook)
			assert.Equal(t, handling.ErrNotPresented, err)
		})
	}
}

func TestFactory_GetHook_UnknownHook(t *testing.T) {
	factory := handling.NewFactory(
		mocks.NewEngineMock(t),
		configuration.HooksConfig{},
	)

	hook, err := factory.GetHook("unknown-hook")

	assert.Nil(t, hook)
	assert.EqualError(t, err, "unknown hook")
}

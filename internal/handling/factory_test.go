package handling_test

import (
	"fisherman/internal/configuration"
	"fisherman/internal/constants"
	. "fisherman/internal/handling"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetHook(t *testing.T) {
	var globalVars = map[string]interface{}{}

	t.Run("return correct hook", func(t *testing.T) {
		factory := NewHookHandlerFactory(
			mocks.NewEngineMock(t).EvalMock.Return(false, nil),
			configuration.HooksConfig{
				ApplyPatchMsgHook:     &configuration.HookConfig{},
				FsMonitorWatchmanHook: &configuration.HookConfig{},
				PostUpdateHook:        &configuration.HookConfig{},
				PreApplyPatchHook:     &configuration.HookConfig{},
				PreCommitHook:         &configuration.HookConfig{},
				PrePushHook:           &configuration.HookConfig{},
				PreRebaseHook:         &configuration.HookConfig{},
				PreReceiveHook:        &configuration.HookConfig{},
				UpdateHook:            &configuration.HookConfig{},
				CommitMsgHook:         &configuration.HookConfig{},
				PrepareCommitMsgHook:  &configuration.HookConfig{},
			},
		)

		for _, tt := range constants.HooksNames {
			t.Run(tt, func(t *testing.T) {
				hook, err := factory.GetHook(tt, globalVars)

				assert.NotNil(t, hook)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("not configured hooks", func(t *testing.T) {
		factory := NewHookHandlerFactory(
			mocks.NewEngineMock(t),
			configuration.HooksConfig{},
		)

		for _, tt := range constants.HooksNames {
			t.Run(tt, func(t *testing.T) {
				hook, err := factory.GetHook(tt, globalVars)

				assert.Nil(t, hook)
				assert.Equal(t, ErrNotPresented, err)
			})
		}
	})

	t.Run("unknown hook", func(t *testing.T) {
		factory := NewHookHandlerFactory(
			mocks.NewEngineMock(t),
			configuration.HooksConfig{},
		)

		hook, err := factory.GetHook("unknown-hook", globalVars)

		assert.Nil(t, hook)
		assert.EqualError(t, err, "unknown hook")
	})
}

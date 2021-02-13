package hookfactory_test

import (
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/hookfactory"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetHook(t *testing.T) {
	factory := hookfactory.NewFactory(
		mocks.NewEngineMock(t),
		configuration.HooksConfig{
			PreCommitHook:        &configuration.PreCommitHookConfig{},
			PrePushHook:          &configuration.PrePushHookConfig{},
			CommitMsgHook:        &configuration.CommitMsgHookConfig{MessagePrefix: "test"},
			PrepareCommitMsgHook: &configuration.PrepareCommitMsgHookConfig{Message: "test"},
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

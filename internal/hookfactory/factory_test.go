package hookfactory_test

import (
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/internal/hookfactory"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactory_GetHook(t *testing.T) {
	shell := configuration.ScriptsConfig{
		"demo": configuration.ScriptConfig{},
	}
	factory := hookfactory.NewFactory(
		mocks.NewCtxFactoryMock(t),
		mocks.NewExtractorMock(t).VariablesMock.Return(map[string]interface{}{}, nil),
		configuration.HooksConfig{
			PreCommitHook:        configuration.PreCommitHookConfig{Shell: shell},
			PrePushHook:          configuration.PrePushHookConfig{Shell: shell},
			CommitMsgHook:        configuration.CommitMsgHookConfig{MessagePrefix: "test"},
			PrepareCommitMsgHook: configuration.PrepareCommitMsgHookConfig{Message: "test"},
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
		mocks.NewCtxFactoryMock(t),
		mocks.NewExtractorMock(t).VariablesMock.Return(map[string]interface{}{}, nil),
		configuration.HooksConfig{},
	)

	for _, tt := range constants.HooksNames {
		t.Run(tt, func(t *testing.T) {
			hook, err := factory.GetHook(tt)

			assert.Nil(t, hook)
			assert.NoError(t, err)
		})
	}
}

func TestFactory_GetHook_UnknownHook(t *testing.T) {
	factory := hookfactory.NewFactory(
		mocks.NewCtxFactoryMock(t),
		mocks.NewExtractorMock(t).VariablesMock.Return(map[string]interface{}{}, nil),
		configuration.HooksConfig{},
	)

	hook, err := factory.GetHook("unknown-hook")

	assert.Nil(t, hook)
	assert.EqualError(t, err, "unknown hook")
}

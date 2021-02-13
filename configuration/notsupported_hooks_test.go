package configuration_test

import (
	"fisherman/configuration"
	"fisherman/internal/hookfactory"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	name   string
	config hookfactory.CompilableConfig
}

var tests []testData = []testData{
	{
		name:   "ApplyPatchMsgHookConfig",
		config: &configuration.ApplyPatchMsgHookConfig{},
	},
	{
		name:   "FsMonitorWatchmanHookConfig",
		config: &configuration.FsMonitorWatchmanHookConfig{},
	},
	{
		name:   "ApplyPatchMsgHookConfig",
		config: &configuration.ApplyPatchMsgHookConfig{},
	},
	{
		name:   "PostUpdateHookConfig",
		config: &configuration.PostUpdateHookConfig{},
	},
	{
		name:   "PreApplyPatchHookConfig",
		config: &configuration.FsMonitorWatchmanHookConfig{},
	},
	{
		name:   "PreRebaseHookConfig",
		config: &configuration.PreRebaseHookConfig{},
	},
	{
		name:   "PreReceiveHookConfig",
		config: &configuration.PreReceiveHookConfig{},
	},
	{
		name:   "UpdateHookConfig",
		config: &configuration.UpdateHookConfig{},
	},
}

func TestFsMonitorWatchmanHookConfig_GetVariablesConfig(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				_ = tt.config.GetVariables()
			})
		})
	}
}

func TestFsMonitorWatchmanHookConfig_Compile(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				tt.config.Compile(mocks.NewEngineMock(t), map[string]interface{}{})
			})
		})
	}
}

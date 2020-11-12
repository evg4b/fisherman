package handlers_test

import (
	"context"
	"errors"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/handlers"
	"fisherman/infrastructure"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrePushHandler(t *testing.T) {
	assert.NotPanics(t, func() {
		ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: mocks.NewRepositoryMock(t).
				GetCurrentBranchMock.Return("develop", nil).
				GetLastTagMock.Return("0.0.0", nil).
				GetUserMock.Return(infrastructure.User{}, nil),
			Shell: mocks.NewShellMock(t),
			App:   &clicontext.AppInfo{},
		})
		err := new(handlers.PrePushHandler).Handle(ctx, []string{})
		assert.NoError(t, err)
	})
}

func TestPrePushHandler_VariablesError(t *testing.T) {
	assert.NotPanics(t, func() {
		ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: mocks.NewRepositoryMock(t).
				GetCurrentBranchMock.Return("develop", nil).
				GetLastTagMock.Return("", errors.New("fail")).
				GetUserMock.Return(infrastructure.User{}, nil),
			Shell: mocks.NewShellMock(t),
			App:   &clicontext.AppInfo{},
		})
		err := new(handlers.PrePushHandler).Handle(ctx, []string{})
		assert.Error(t, err, "fail")
	})
}

func TestPrePushHandler_IsConfigured(t *testing.T) {
	var handler handlers.PrePushHandler

	tests := []struct {
		name     string
		config   *config.HooksConfig
		expected bool
	}{
		{
			name:     "empty structure",
			config:   &config.HooksConfig{},
			expected: false,
		},
		{
			name: "configured script",
			config: &config.HooksConfig{
				PrePushHook: hooks.PrePushHookConfig{
					Shell: hooks.ScriptsConfig{
						"demo": hooks.ScriptConfig{
							Commands: []string{"ls"},
							Output:   true,
						},
					},
					Variables: hooks.Variables{},
				},
			},
			expected: true,
		},
		{
			name: "configured script",
			config: &config.HooksConfig{
				PrePushHook: hooks.PrePushHookConfig{
					Variables: hooks.Variables{
						FromBranch:  "demo",
						FromLastTag: "demo",
					},
				},
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := handler.IsConfigured(tt.config)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

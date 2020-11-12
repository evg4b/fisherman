package handlers_test

import (
	"context"
	"errors"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/handlers"
	"fisherman/infrastructure"
	inf_mock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrePushHandler(t *testing.T) {
	fakeRepository := inf_mock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)
	fakeRepository.On("GetUser").Return(infrastructure.User{}, nil)

	fakeShell := inf_mock.Shell{}

	assert.NotPanics(t, func() {
		ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: &fakeRepository,
			Shell:      &fakeShell,
			App:        &clicontext.AppInfo{},
		})
		err := new(handlers.PrePushHandler).Handle(ctx, []string{})
		assert.NoError(t, err)
	})
}

func TestPrePushHandler_VariablesError(t *testing.T) {
	fakeRepository := inf_mock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", errors.New("fail"))
	fakeRepository.On("GetUser").Return(infrastructure.User{}, nil)

	fakeShell := inf_mock.Shell{}

	assert.NotPanics(t, func() {
		ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: &fakeRepository,
			Shell:      &fakeShell,
			App:        &clicontext.AppInfo{},
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

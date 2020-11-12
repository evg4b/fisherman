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

func TestPreCommitHandler(t *testing.T) {
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
		err := new(handlers.PreCommitHandler).Handle(ctx, []string{})
		assert.NoError(t, err)
	})
}

func TestPreCommitHandler_VariablesError(t *testing.T) {
	var handler handlers.PreCommitHandler

	assert.NotPanics(t, func() {
		ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
			GlobalVariables: map[string]interface{}{},
			Config: &config.FishermanConfig{
				Hooks: config.HooksConfig{},
			},
			Repository: mocks.NewRepositoryMock(t).
				GetCurrentBranchMock.Return("", errors.New("fail")).
				GetLastTagMock.Return("0.0.0", nil).
				GetUserMock.Return(infrastructure.User{}, nil),
			Shell: mocks.NewShellMock(t),
			App:   &clicontext.AppInfo{},
		})
		err := handler.Handle(ctx, []string{})
		assert.Error(t, err, "fail")
	})
}

func TestPreCommitHandler_IsConfigured(t *testing.T) {
	var handler handlers.PreCommitHandler

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
			name: "empty shell scripts structure",
			config: &config.HooksConfig{
				PreCommitHook: hooks.PreCommitHookConfig{
					Shell: make(hooks.ScriptsConfig),
				},
			},
			expected: false,
		},
		{
			name: "not empty shell scripts structure",
			config: &config.HooksConfig{
				PreCommitHook: hooks.PreCommitHookConfig{
					Shell: hooks.ScriptsConfig{
						"demo": hooks.ScriptConfig{},
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

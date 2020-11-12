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

func TestPreCommitHandler(t *testing.T) {
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
		err := new(handlers.PreCommitHandler).Handle(ctx, []string{})
		assert.NoError(t, err)
	})
}

func TestPreCommitHandler_VariablesError(t *testing.T) {
	fakeRepository := inf_mock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("", errors.New("fail"))
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)
	fakeRepository.On("GetUser").Return(infrastructure.User{}, nil)

	fakeShell := inf_mock.Shell{}

	var handler handlers.PreCommitHandler

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

package handlers

import (
	"context"
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/config/hooks"
	inf_mock "fisherman/mocks/infrastructure"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCommitMsgHandler(t *testing.T) {
	fakeRepository := inf_mock.Repository{}
	fakeRepository.On("GetCurrentBranch").Return("develop", nil)
	fakeRepository.On("GetLastTag").Return("0.0.0", nil)

	fakeFS := inf_mock.FileSystem{}
	fakeFS.On("Read", ".git/MESSAGE").Return("[fisherman] test commit", nil)

	tests := []struct {
		name string
		args []string
		err  error
	}{
		{name: "base test", args: []string{}, err: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
				Config:     &config.DefaultConfig,
				Repository: &fakeRepository,
				FileSystem: &fakeFS,
			})
			err := new(PrepareCommitMsgHandler).Handle(ctx, tt.args)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestPrepareCommitMsgHandler_IsConfigured(t *testing.T) {
	var handler PrepareCommitMsgHandler
	tests := []struct {
		name     string
		config   *config.HooksConfig
		expected bool
	}{
		{name: "empty structure", config: &config.HooksConfig{}, expected: false},
		{
			name: "empty PrepareCommitMsgHookConfig structure",
			config: &config.HooksConfig{
				PrepareCommitMsgHook: hooks.PrepareCommitMsgHookConfig{},
			},
			expected: false,
		},
		{
			name: "empty PrepareCommitMsgHookConfig structure",
			config: &config.HooksConfig{
				PrepareCommitMsgHook: hooks.PrepareCommitMsgHookConfig{
					Message: "demo",
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

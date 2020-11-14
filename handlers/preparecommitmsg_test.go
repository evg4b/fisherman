package handlers

import (
	"context"
	"fisherman/config"
	"fisherman/config/hooks"
	"fisherman/internal/clicontext"
	"fisherman/mocks"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCommitMsgHandler(t *testing.T) {
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
				Config: &config.DefaultConfig,
				Repository: mocks.NewRepositoryMock(t).
					GetCurrentBranchMock.Return("develop", nil).
					GetLastTagMock.Return("0.0.0", nil),
				FileSystem: mocks.NewFileSystemMock(t).
					ReadMock.When(".git/MESSAGE").Then("[fisherman] test commit", nil),
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

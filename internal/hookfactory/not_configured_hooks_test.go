package hookfactory_test

import (
	"context"
	c "fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
	"fisherman/internal/hookfactory"
	"fisherman/mocks"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguredHooks(t *testing.T) {
	ctxFactory := func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(
			context.TODO(),
			mocks.NewFileSystemMock(t),
			mocks.NewShellMock(t),
			mocks.NewRepositoryMock(t),
			args,
			output,
		)
	}

	vars := c.VariablesConfig{FromBranch: "Data"}

	tests := []testCase{
		{
			name: "pre-push",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PrePush(c.PrePushHookConfig{Variables: vars})
			},
		},
		{
			name: "pre-commit",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PreCommit(c.PreCommitHookConfig{Variables: vars})
			},
		},
		{
			name: "commit-msg",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.CommitMsg(c.CommitMsgHookConfig{NotEmpty: true})
			},
		},
		{
			name: "prepare-commit-msg",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PrepareCommitMsg(c.PrepareCommitMsgHookConfig{Variables: vars})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled := false
			factory := hookfactory.NewFactory(ctxFactory, func(configcompiler.CompilableConfig) {
				compiled = true
			})
			result := tt.constructor(factory)

			assert.IsType(t, &handling.HookHandler{}, result.Handler)
			assert.True(t, result.Registered)
			assert.True(t, compiled)
		})
	}
}

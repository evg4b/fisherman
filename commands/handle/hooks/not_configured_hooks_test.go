package hooks_test

import (
	"context"
	"fisherman/commands/handle/hooks"
	c "fisherman/config/hooks"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/handling"
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

	sh := mocks.NewShellMock(t)
	vars := c.Variables{FromBranch: "Data"}

	tests := []struct {
		name        string
		constructor constructorFunc
	}{
		{
			name: "pre-push",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PrePush(ctxFactory, c.PrePushHookConfig{Variables: vars}, sh, compile)
			},
		},
		{
			name: "pre-commit",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PreCommit(ctxFactory, c.PreCommitHookConfig{Variables: vars}, sh, compile)
			},
		},
		{
			name: "commit-msg",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.CommitMsg(ctxFactory, c.CommitMsgHookConfig{NotEmpty: true}, compile)
			},
		},
		{
			name: "prepare-commit-msg",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PrepareCommitMsg(ctxFactory, c.PrepareCommitMsgHookConfig{Variables: vars}, compile)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled := false
			result := tt.constructor(ctxFactory, func(configcompiler.CompilableConfig) {
				compiled = true
			})

			assert.IsType(t, &handling.HookHandler{}, result.Handler)
			assert.True(t, result.Registered)
			assert.True(t, compiled)
		})
	}
}

package hooks_test

import (
	"context"
	"fisherman/commands/handle/hooks"
	c "fisherman/config/hooks"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/mocks"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type constructorFunc = func(internal.CtxFactory, configcompiler.Compiler) hooks.HandlerRegistration

func TestNotConfiguredHooks(t *testing.T) {
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

	compiler := func(configcompiler.CompilableConfig) {
		t.Fatal("compile should not be called")
	}

	tests := []struct {
		name        string
		constructor constructorFunc
	}{
		{
			name: "pre-push",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PrePush(ctxFactory, c.PrePushHookConfig{}, sh, compile)
			},
		},
		{
			name: "pre-commit",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PreCommit(ctxFactory, c.PreCommitHookConfig{}, sh, compile)
			},
		},
		{
			name: "commit-msg",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.CommitMsg(ctxFactory, c.CommitMsgHookConfig{}, compile)
			},
		},
		{
			name: "prepare-commit-msg",
			constructor: func(ctxFactory internal.CtxFactory, compile configcompiler.Compiler) hooks.HandlerRegistration {
				return hooks.PrepareCommitMsg(ctxFactory, c.PrepareCommitMsgHookConfig{}, compile)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.constructor(ctxFactory, compiler)

			assert.ObjectsAreEqual(hooks.NotRegistered, result)
		})
	}
}

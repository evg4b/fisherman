package hookfactory_test

import (
	"context"
	c "fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/configcompiler"
	"fisherman/internal/hookfactory"
	"fisherman/mocks"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type constructorFunc = func(*hookfactory.Factory) hookfactory.HandlerRegistration
type testCase struct {
	name        string
	constructor constructorFunc
}

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

	compiler := func(configcompiler.CompilableConfig) {
		t.Fatal("compile should not be called")
	}

	tests := []testCase{
		{
			name: "pre-push",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PrePush(c.PrePushHookConfig{})
			},
		},
		{
			name: "pre-commit",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PreCommit(c.PreCommitHookConfig{})
			},
		},
		{
			name: "commit-msg",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.CommitMsg(c.CommitMsgHookConfig{})
			},
		},
		{
			name: "prepare-commit-msg",
			constructor: func(factory *hookfactory.Factory) hookfactory.HandlerRegistration {
				return factory.PrepareCommitMsg(c.PrepareCommitMsgHookConfig{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.constructor(hookfactory.NewFactory(ctxFactory, compiler))

			assert.ObjectsAreEqual(hookfactory.NotRegistered, result)
		})
	}
}

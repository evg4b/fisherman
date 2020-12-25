package mocks

import (
	"context"
	"fisherman/internal"
	"io"
	"testing"
)

func NewCtxFactoryMock(t *testing.T) internal.CtxFactory {
	return func(args []string, output io.Writer) *internal.Context {
		return internal.NewInternalContext(
			context.TODO(),
			NewFileSystemMock(t),
			NewShellMock(t),
			NewRepositoryMock(t),
			args,
			output,
		)
	}
}

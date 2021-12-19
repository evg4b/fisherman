package handle_test

import (
	. "fisherman/internal/commands/handle"
	"fisherman/testing/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := NewCommand(
		WithFileSystem(mocks.NewFilesystemMock(t)),
		WithRepository(mocks.NewRepositoryMock(t)),
		WithCwd("/"),
		WithExpressionEngine(mocks.NewEngineMock(t)),
		WithHooksConfig(&mocks.HooksConfigStub),
	)

	assert.Equal(t, "handle", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand(
		WithFileSystem(mocks.NewFilesystemMock(t)),
		WithRepository(mocks.NewRepositoryMock(t)),
		WithCwd("/"),
		WithExpressionEngine(mocks.NewEngineMock(t)),
		WithHooksConfig(&mocks.HooksConfigStub),
	)

	assert.NotEmpty(t, command.Description())
}

func TestNewCommand(t *testing.T) {
	t.Run("panic when cwd is not configured", func(t *testing.T) {
		t.SkipNow()
		assert.PanicsWithError(t, "Cwd should be configured", func() {
			_ = NewCommand(
				WithFileSystem(mocks.NewFilesystemMock(t)),
				WithRepository(mocks.NewRepositoryMock(t)),
				WithExpressionEngine(mocks.NewEngineMock(t)),
				WithHooksConfig(&mocks.HooksConfigStub),
			)
		})
	})

	t.Run("panic when filesystem is not configured", func(t *testing.T) {
		t.SkipNow()
		assert.PanicsWithError(t, "FileSystem should be configured", func() {
			_ = NewCommand(
				WithRepository(mocks.NewRepositoryMock(t)),
				WithCwd("/"),
				WithExpressionEngine(mocks.NewEngineMock(t)),
				WithHooksConfig(&mocks.HooksConfigStub),
			)
		})
	})

	t.Run("panic when repository is not configured", func(t *testing.T) {
		t.SkipNow()
		assert.PanicsWithError(t, "Repository should be configured", func() {
			_ = NewCommand(
				WithFileSystem(mocks.NewFilesystemMock(t)),
				WithCwd("/"),
				WithExpressionEngine(mocks.NewEngineMock(t)),
				WithHooksConfig(&mocks.HooksConfigStub),
			)
		})
	})

	t.Run("panic when expression engine is not configured", func(t *testing.T) {
		t.SkipNow()
		assert.PanicsWithError(t, "ExpressionEngine should be configured", func() {
			_ = NewCommand(
				WithFileSystem(mocks.NewFilesystemMock(t)),
				WithRepository(mocks.NewRepositoryMock(t)),
				WithCwd("/"),
				WithHooksConfig(&mocks.HooksConfigStub),
			)
		})
	})

	t.Run("panic when hooks config is not configured", func(t *testing.T) {
		assert.PanicsWithError(t, "HooksConfig should be configured", func() {
			_ = NewCommand(
				WithFileSystem(mocks.NewFilesystemMock(t)),
				WithRepository(mocks.NewRepositoryMock(t)),
				WithCwd("/"),
				WithExpressionEngine(mocks.NewEngineMock(t)),
			)
		})
	})
}

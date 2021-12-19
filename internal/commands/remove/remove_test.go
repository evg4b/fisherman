package remove_test

import (
	"context"
	"errors"
	. "fisherman/internal/commands/remove"
	"fisherman/internal/configuration"
	"fisherman/pkg/log"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/stretchr/testify/assert"
)

// nolint: gochecknoinits
func init() {
	log.SetOutput(io.Discard)
}

func TestCommand_Run(t *testing.T) {
	cwd := filepath.Join("usr", "home")
	configs := map[string]string{
		configuration.GlobalMode: filepath.Join("usr", "home", ".fisherman.yml"),
	}

	tests := []struct {
		name        string
		fs          billy.Filesystem
		expectedErr string
	}{
		{
			name: "executed successful",
			fs: testutils.FsFromMap(t, map[string]string{
				filepath.Join("usr", "home", ".fisherman.yml"):              "content",
				filepath.Join("usr", "home", ".fisherman.yaml"):             "content",
				filepath.Join("usr", "home", ".git", "hooks", "commit-msg"): "content",
				filepath.Join("usr", "home", ".fisherman.yml"):              "content",
			}),
		},
		{
			name: "exist errors",
			fs: mocks.NewFilesystemMock(t).
				StatMock.Return(nil, errors.New("test error")),
			expectedErr: "test error",
		},
		{
			name: "delete error",
			fs: mocks.NewFilesystemMock(t).
				StatMock.Return(nil, nil).
				RemoveMock.Return(errors.New("delete error")),
			expectedErr: "delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := NewCommand(
				WithFileSystem(tt.fs),
				WithCwd(cwd),
				WithConfigFiles(configs),
			)

			err := command.Run(context.TODO(), []string{})

			testutils.AssertError(t, tt.expectedErr, err)
		})
	}
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand(
		WithFileSystem(mocks.NewFilesystemMock(t)),
		WithCwd("/"),
	)

	assert.Equal(t, "remove", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand(
		WithFileSystem(mocks.NewFilesystemMock(t)),
		WithCwd("/"),
	)

	assert.NotEmpty(t, command.Description())
}

func TestNewCommand(t *testing.T) {
	t.Run("panic when cwd is not configured", func(t *testing.T) {
		fs := mocks.NewFilesystemMock(t)
		assert.PanicsWithError(t, "Cwd should be configured", func() {
			_ = NewCommand(WithFileSystem(fs))
		})
	})

	t.Run("panic when filesystem is not configured", func(t *testing.T) {
		assert.PanicsWithError(t, "FileSystem should be configured", func() {
			_ = NewCommand(WithCwd("/"))
		})
	})
}

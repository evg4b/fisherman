package remove_test

import (
	"errors"
	"fisherman/internal"
	. "fisherman/internal/commands/remove"
	"fisherman/internal/configuration"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	appInfo := internal.AppInfo{
		Cwd:        filepath.Join("usr", "home"),
		Executable: filepath.Join("bin", "fisherman.exe"),
		Configs: map[string]string{
			configuration.GlobalMode: filepath.Join("usr", "home", ".fisherman.yml"),
		},
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
			command := NewCommand(tt.fs, appInfo, &user.User{
				HomeDir: filepath.Join("usr", "home"),
			})

			err := command.Init([]string{})
			assert.NoError(t, err)

			err = command.Run(mocks.NewExecutionContextMock(t))

			testutils.AssertError(t, tt.expectedErr, err)
		})
	}
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand(
		mocks.NewFilesystemMock(t),
		mocks.AppInfoStub,
		&testutils.TestUser,
	)

	assert.Equal(t, "remove", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand(
		mocks.NewFilesystemMock(t),
		mocks.AppInfoStub,
		&testutils.TestUser,
	)

	assert.NotEmpty(t, command.Description())
}

package remove_test

import (
	"errors"
	"fisherman/commands/remove"
	"fisherman/configuration"
	"fisherman/infrastructure"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCommand_Run(t *testing.T) {
	fs := testutils.FsFromMap(t, map[string]string{
		filepath.Join("usr", "home", ".fisherman.yml"):              "content",
		filepath.Join("usr", "home", ".fisherman.yaml"):             "content",
		filepath.Join("usr", "home", ".git", "hooks", "commit-msg"): "content",
		filepath.Join("usr", "home", ".fisherman.yml"):              "content",
	})

	command := remove.NewCommand(
		fs,
		&internal.AppInfo{
			Cwd:        filepath.Join("usr", "home"),
			Executable: filepath.Join("bin", "fisherman.exe"),
		},
		&user.User{
			HomeDir: filepath.Join("usr", "home"),
		},
	)
	err := command.Init([]string{})
	assert.NoError(t, err)

	err = command.Run()
	assert.NoError(t, err)
}

func TestCommand_Run_WithError(t *testing.T) {
	appInfo := internal.AppInfo{
		Cwd:        filepath.Join("usr", "home"),
		Executable: filepath.Join("bin", "fisherman.exe"),
		Configs: map[string]string{
			configuration.GlobalMode: filepath.Join("usr", "home", ".fisherman.yml"),
		},
	}

	tests := []struct {
		name          string
		files         infrastructure.FileSystem
		expectedError string
	}{
		{
			name:          "exist errors",
			files:         mocks.NewFileSystemMock(t).StatMock.Return(nil, errors.New("Test error")),
			expectedError: "Test error",
		},
		{
			name: "delete error",
			files: mocks.NewFileSystemMock(t).
				StatMock.Return(nil, nil).
				RemoveMock.Return(errors.New("delete error")),
			expectedError: "delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := remove.NewCommand(tt.files, &appInfo, &user.User{
				HomeDir: filepath.Join("usr", "home"),
			})

			err := c.Init([]string{})
			assert.NoError(t, err)

			err = c.Run()

			testutils.CheckError(t, tt.expectedError, err)
		})
	}
}

func TestCommand_Name(t *testing.T) {
	command := remove.NewCommand(
		mocks.NewFileSystemMock(t),
		&internal.AppInfo{},
		&user.User{},
	)

	assert.Equal(t, "remove", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := remove.NewCommand(
		mocks.NewFileSystemMock(t),
		&internal.AppInfo{},
		&user.User{},
	)

	assert.NotEmpty(t, command.Description())
}

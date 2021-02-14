package remove_test

import (
	"errors"
	"fisherman/commands/remove"
	"fisherman/configuration"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/testing/mocks"
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
	command := remove.NewCommand(
		makeFakeFS(t).
			DeleteMock.When(filepath.Join("usr", "home", ".fisherman.yml")).Then(nil).
			ExistMock.When(filepath.Join("usr", "home", ".fisherman.yml")).Then(true).
			ExistMock.When(filepath.Join("usr", "home", ".fisherman.yaml")).Then(true),
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
	expectedError := errors.New("Test error")
	c := remove.NewCommand(
		makeFakeFS(t).DeleteMock.Expect(filepath.Join("usr", "home", ".fisherman.yml")).Return(expectedError),
		&internal.AppInfo{
			Cwd:        filepath.Join("usr", "home"),
			Executable: filepath.Join("bin", "fisherman.exe"),
			Configs: map[string]string{
				configuration.GlobalMode: filepath.Join("usr", "home", ".fisherman.yml"),
			},
		},
		&user.User{
			HomeDir: filepath.Join("usr", "home"),
		},
	)
	err := c.Init([]string{})
	assert.NoError(t, err)

	err = c.Run()
	assert.Equal(t, err, expectedError)
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

func makeFakeFS(t *testing.T) *mocks.FileSystemMock {
	mock := mocks.NewFileSystemMock(t)
	for _, name := range constants.HooksNames {
		mock.ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", name)).Then(false)
	}

	return mock.
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Then(true).
		ExistMock.When(filepath.Join("usr", "home", ".fisherman.yml")).Then(true).
		DeleteMock.Expect(filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(nil)
}

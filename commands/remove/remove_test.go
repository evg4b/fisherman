package remove_test

import (
	"context"
	"errors"
	"fisherman/clicontext"
	"fisherman/commands/remove"
	"fisherman/config"
	"fisherman/mocks"
	"flag"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	fakeFS := makeFakeFS(t)
	fakeFS.DeleteMock.When(filepath.Join("usr", "home", ".fisherman.yml")).Then(nil)

	ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
		App: &clicontext.AppInfo{
			Cwd:        filepath.Join("usr", "home"),
			Executable: filepath.Join("bin", "fisherman.exe"),
		},
		FileSystem: fakeFS,
		Config:     &config.DefaultConfig,
		User: &user.User{
			HomeDir: filepath.Join("usr", "home"),
		},
	})

	c := remove.NewCommand(flag.ExitOnError)
	err := c.Init([]string{})
	assert.NoError(t, err)

	err = c.Run(ctx)
	assert.NoError(t, err)
}

func TestCommand_Run_WithError(t *testing.T) {
	expectedError := errors.New("Test error")
	fakeFS := makeFakeFS(t)
	fakeFS.DeleteMock.Expect(filepath.Join("usr", "home", ".fisherman.yml")).Return(expectedError)

	ctx := clicontext.NewContext(context.TODO(), clicontext.Args{
		App: &clicontext.AppInfo{
			Cwd:        filepath.Join("usr", "home"),
			Executable: filepath.Join("bin", "fisherman.exe"),
		},
		FileSystem: fakeFS,
		Config:     &config.DefaultConfig,
		User: &user.User{
			HomeDir: filepath.Join("usr", "home"),
		},
	})

	c := remove.NewCommand(flag.ExitOnError)
	err := c.Init([]string{})
	assert.NoError(t, err)

	err = c.Run(ctx)
	assert.Equal(t, err, expectedError)
}

func TestCommand_Name(t *testing.T) {
	c := remove.NewCommand(flag.ExitOnError)
	assert.Equal(t, "remove", c.Name())
}

func makeFakeFS(t *testing.T) *mocks.FileSystemMock {
	return mocks.NewFileSystemMock(t).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Then(true).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "prepare-commit-msg")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "pre-commit")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".git", ".fisherman.yml")).Then(false).
		ExistMock.When(filepath.Join("usr", "home", ".fisherman.yml")).Then(true).
		DeleteMock.Expect(filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(nil)
}

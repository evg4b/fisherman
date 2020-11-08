package remove_test

import (
	"errors"
	"fisherman/clicontext"
	"fisherman/commands/remove"
	"fisherman/config"
	inf_mocks "fisherman/mocks/infrastructure"
	"flag"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	fakeFS := makeFakeFS()
	fakeFS.On("Delete", filepath.Join("usr", "home", ".fisherman.yml")).Return(nil)

	ctx := clicontext.NewContext(clicontext.Args{
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
	fakeFS := makeFakeFS()
	fakeFS.On("Delete", filepath.Join("usr", "home", ".fisherman.yml")).Return(expectedError)

	ctx := clicontext.NewContext(clicontext.Args{
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

func makeFakeFS() *inf_mocks.FileSystem {
	fakeFS := inf_mocks.FileSystem{}

	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(true)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "prepare-commit-msg")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-commit")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".git", ".fisherman.yml")).Return(false)
	fakeFS.On("Exist", filepath.Join("usr", "home", ".fisherman.yml")).Return(true)

	fakeFS.On("Delete", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(nil)

	return &fakeFS
}

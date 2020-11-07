package remove_test

import (
	"fisherman/clicontext"
	"fisherman/commands/remove"
	"fisherman/config"
	infmocks "fisherman/mocks/infrastructure"
	"flag"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	fakeFileAccessor := infmocks.FileAccessor{}

	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(true)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "prepare-commit-msg")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-commit")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", ".fisherman.yml")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".fisherman.yml")).Return(true)

	fakeFileAccessor.On("Delete", filepath.Join("usr", "home", ".fisherman.yml")).Return(nil)
	fakeFileAccessor.On("Delete", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(nil)

	tests := []struct {
		name     string
		expected error
	}{
		{name: "demo", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := clicontext.NewContext(clicontext.Args{
				App: &clicontext.AppInfo{
					Cwd:        filepath.Join("usr", "home"),
					Executable: filepath.Join("bin", "fisherman.exe"),
				},
				FileAccessor: &fakeFileAccessor,
				Config:       &config.DefaultConfig,
				User: &user.User{
					HomeDir: filepath.Join("usr", "home"),
				},
			})

			c := remove.NewCommand(flag.ExitOnError)
			err := c.Init([]string{})
			assert.NoError(t, err)
			err = c.Run(ctx)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestCommand_Name(t *testing.T) {
	c := remove.NewCommand(flag.ExitOnError)
	assert.Equal(t, "remove", c.Name())
}

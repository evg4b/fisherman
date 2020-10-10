package remove_test

import (
	"fisherman/commands"
	"fisherman/commands/remove"
	"fisherman/config"
	infmocks "fisherman/mocks/infrastructure"
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	fakeFileAccessor := infmocks.FileAccessor{}

	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(true)
	fakeFileAccessor.On("Delete", filepath.Join("usr", "home", ".git", "hooks", "commit-msg")).Return(nil)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "prepare-commit-msg")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-commit")).Return(false)
	fakeFileAccessor.On("Exist", filepath.Join("usr", "home", ".git", "hooks", "pre-push")).Return(false)

	tests := []struct {
		name     string
		expected error
	}{
		{name: "demo", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := commands.NewContext(commands.CliCommandContextParams{
				App: &commands.AppInfo{
					Cwd:        filepath.Join("usr", "home"),
					Executable: filepath.Join("bin", "fisherman.exe"),
				},
				FileAccessor: &fakeFileAccessor,
				Config:       &config.DefaultConfig,
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

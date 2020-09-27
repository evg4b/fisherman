package init_test

import (
	"fisherman/commands"
	initc "fisherman/commands/init"
	"fisherman/config"
	"fisherman/constants"
	iomock "fisherman/mocks/infrastructure/io"
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCommand(t *testing.T) {
	command := initc.NewCommand(flag.ExitOnError)
	assert.NotNil(t, command)
}

func TestCommand_Run_Force_Mode(t *testing.T) {
	command := initc.NewCommand(flag.ExitOnError)
	cwd := "/demo/"

	fakeFileAccessor := iomock.FileAccessor{}
	fakeFileAccessor.On("Write", mock.IsType("string"), mock.IsType("string")).Return(nil)
	fakeFileAccessor.On("Exist", filepath.Join(cwd, constants.AppConfigName)).Return(true)

	tests := []struct {
		name string
		args []string
	}{
		{name: "dem", args: []string{"--force"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := commands.CommandContext{
				Files: &fakeFileAccessor,
				App: &commands.AppInfo{
					Cwd:                cwd,
					IsRegisteredInPath: true,
				},
				Config: &config.DefaultConfig,
			}
			command.Init(tt.args)
			err := command.Run(&ctx)
			assert.NoError(t, err)
		})
	}
}

func TestCommand_Name(t *testing.T) {
	command := initc.NewCommand(flag.ExitOnError)
	assert.Equal(t, command.Name(), "init")
}

package initialize_test

import (
	"fisherman/commands/initialize"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/internal/clicontext"
	"fisherman/mocks"
	"flag"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	command := initialize.NewCommand(flag.ExitOnError)
	assert.NotNil(t, command)
}

func TestCommand_Run_Force_Mode(t *testing.T) {
	command := initialize.NewCommand(flag.ExitOnError)
	cwd := "/demo/"

	d := mocks.NewFileSystemMock(t).
		WriteMock.Return(nil).
		ExistMock.When(filepath.Join(cwd, constants.AppConfigName)).Then(true).
		ChmodMock.Return(nil).
		ChownMock.Return(nil)

	ctx := clicontext.CommandContext{
		Files: d,
		App: &clicontext.AppInfo{
			Cwd: cwd,
		},
		Config: &config.HooksConfig{},
		User:   &user.User{},
	}
	err := command.Init([]string{"--force"})
	assert.NoError(t, err)
	err = command.Run(&ctx)
	assert.NoError(t, err)
}

func TestCommand_Name(t *testing.T) {
	command := initialize.NewCommand(flag.ExitOnError)
	assert.Equal(t, command.Name(), "init")
}

package initialize_test

import (
	"fisherman/clicontext"
	"fisherman/commands/initialize"
	"fisherman/config"
	"fisherman/constants"
	inf_mock "fisherman/mocks/infrastructure"
	"flag"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCommand(t *testing.T) {
	command := initialize.NewCommand(flag.ExitOnError)
	assert.NotNil(t, command)
}

func TestCommand_Run_Force_Mode(t *testing.T) {
	user := user.User{}
	command := initialize.NewCommand(flag.ExitOnError)
	cwd := "/demo/"

	fakeFS := inf_mock.FileSystem{}
	fakeFS.On("Write", mock.IsType("string"), mock.IsType("string")).Return(nil)
	fakeFS.On("Exist", filepath.Join(cwd, constants.AppConfigName)).Return(true)
	fakeFS.On("Chmod", mock.IsType("string"), os.ModePerm).Return(nil)
	fakeFS.On("Chown", mock.IsType("string"), &user).Return(nil)

	ctx := clicontext.CommandContext{
		Files: &fakeFS,
		App: &clicontext.AppInfo{
			Cwd: cwd,
		},
		Config: &config.HooksConfig{},
		User:   &user,
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

package initialize_test

import (
	"fisherman/commands/initialize"
	"fisherman/constants"
	"fisherman/internal"
	"fisherman/mocks"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	command := initialize.NewCommand(mocks.NewFileSystemMock(t), &internal.AppInfo{}, &user.User{})
	assert.NotNil(t, command)
}

func TestCommand_Run_Force_Mode(t *testing.T) {
	cwd := "/demo/"

	command := initialize.NewCommand(
		mocks.NewFileSystemMock(t).
			WriteMock.Return(nil).
			ExistMock.When(filepath.Join(cwd, constants.AppConfigName)).Then(true).
			ChmodMock.Return(nil).
			ChownMock.Return(nil),
		&internal.AppInfo{Cwd: cwd},
		&user.User{},
	)

	err := command.Init([]string{"--force"})
	assert.NoError(t, err)
	err = command.Run()
	assert.NoError(t, err)
}

func TestCommand_Name(t *testing.T) {
	command := initialize.NewCommand(mocks.NewFileSystemMock(t), &internal.AppInfo{}, &user.User{})
	assert.Equal(t, command.Name(), "init")
}

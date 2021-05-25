package initialize_test

import (
	"fisherman/commands/initialize"
	"fisherman/internal"
	"fisherman/internal/constants"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	command := initialize.NewCommand(mocks.NewFileSystemMock(t), mocks.AppInfoStub, &testutils.TestUser)
	assert.NotNil(t, command)
}

func TestCommand_Run_Force_Mode(t *testing.T) {
	cwd := "/demo/"

	fs := testutils.FsFromMap(t, map[string]string{
		filepath.Join(cwd, constants.AppConfigNames[0]): "content",
	})

	command := initialize.NewCommand(fs, internal.AppInfo{Cwd: cwd}, &testutils.TestUser)

	err := command.Init([]string{"--force"})
	assert.NoError(t, err)
	err = command.Run(mocks.NewExecutionContextMock(t))
	assert.NoError(t, err)
}

func TestCommand_Name(t *testing.T) {
	command := initialize.NewCommand(mocks.NewFileSystemMock(t), mocks.AppInfoStub, &testutils.TestUser)

	assert.Equal(t, command.Name(), "init")
}

func TestCommand_Description(t *testing.T) {
	command := initialize.NewCommand(mocks.NewFileSystemMock(t), mocks.AppInfoStub, &testutils.TestUser)

	assert.NotEmpty(t, command.Description())
}

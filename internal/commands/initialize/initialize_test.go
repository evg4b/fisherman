package initialize_test

import (
	"context"
	"fisherman/internal"
	. "fisherman/internal/commands/initialize"
	"fisherman/internal/constants"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	command := NewCommand(mocks.NewFilesystemMock(t), mocks.AppInfoStub, &testutils.TestUser)
	assert.NotNil(t, command)
}

func TestCommand_Run(t *testing.T) {
	t.Run("runs with force mode", func(t *testing.T) {
		cwd := "/demo/"

		fs := testutils.FsFromMap(t, map[string]string{
			filepath.Join(cwd, constants.AppConfigNames[0]): "content",
		})

		command := NewCommand(fs, internal.AppInfo{Cwd: cwd}, &testutils.TestUser)

		err := command.Run(context.TODO(), []string{"--force"})
		assert.NoError(t, err)
	})
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand(mocks.NewFilesystemMock(t), mocks.AppInfoStub, &testutils.TestUser)

	assert.Equal(t, command.Name(), "init")
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand(mocks.NewFilesystemMock(t), mocks.AppInfoStub, &testutils.TestUser)

	assert.NotEmpty(t, command.Description())
}

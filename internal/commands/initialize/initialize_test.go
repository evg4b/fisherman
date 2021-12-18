package initialize_test

import (
	"context"
	. "fisherman/internal/commands/initialize"
	"fisherman/internal/constants"
	"fisherman/pkg/log"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: gochecknoinits
func init() {
	log.SetOutput(io.Discard)
}

func TestNewCommand(t *testing.T) {
	command := NewCommand()

	assert.NotNil(t, command)
}

func TestCommand_Run(t *testing.T) {
	t.Run("runs with force mode", func(t *testing.T) {
		cwd := "/demo/"

		fs := testutils.FsFromMap(t, map[string]string{
			filepath.Join(cwd, constants.AppConfigNames[0]): "content",
		})

		command := NewCommand(WithFilesystem(fs), WithCwd(cwd), WithUser(&testutils.TestUser))

		err := command.Run(context.TODO(), []string{"--force"})
		assert.NoError(t, err)
	})
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand(
		WithFilesystem(mocks.NewFilesystemMock(t)),
		WithCwd(mocks.Cwd),
		WithExecutable(mocks.Executable),
		WithUser(&testutils.TestUser),
	)

	assert.Equal(t, command.Name(), "init")
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand(
		WithFilesystem(mocks.NewFilesystemMock(t)),
		WithCwd(mocks.Cwd),
		WithExecutable(mocks.Executable),
		WithUser(&testutils.TestUser),
	)

	assert.NotEmpty(t, command.Description())
}

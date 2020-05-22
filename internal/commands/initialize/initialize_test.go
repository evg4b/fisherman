package initialize_test

import (
	"context"
	"github.com/evg4b/fisherman/internal/constants"
	"github.com/evg4b/fisherman/pkg/log"
	"github.com/evg4b/fisherman/testing/mocks"
	"github.com/evg4b/fisherman/testing/testutils"
	"io"
	"path/filepath"
	"testing"

	. "github.com/evg4b/fisherman/internal/commands/initialize"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint: gochecknoinits
func init() {
	log.SetOutput(io.Discard)
}

func TestNewCommand(t *testing.T) {
	command := NewCommand(
		WithFilesystem(mocks.NewFilesystemMock(t)),
		WithCwd(mocks.Cwd),
		WithExecutable(mocks.Executable),
		WithUser(&testutils.TestUser),
	)

	assert.NotNil(t, command)
}

func TestCommand_Run(t *testing.T) {
	t.Run("runs with force mode", func(t *testing.T) {
		cwd := "/demo/"

		fs := testutils.FsFromMap(t, map[string]string{
			filepath.Join(cwd, constants.AppConfigNames[0]): "content",
		})

		command := NewCommand(
			WithFilesystem(fs),
			WithCwd(cwd),
			WithExecutable(mocks.Executable),
			WithUser(&testutils.TestUser),
		)

		err := command.Run(context.TODO(), []string{"--force"})
		require.NoError(t, err)
	})
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand(
		WithFilesystem(mocks.NewFilesystemMock(t)),
		WithCwd(mocks.Cwd),
		WithExecutable(mocks.Executable),
		WithUser(&testutils.TestUser),
	)

	assert.Equal(t, "init", command.Name())
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

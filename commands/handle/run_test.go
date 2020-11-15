package handle_test

import (
	"errors"
	"fisherman/commands/handle"
	"fisherman/config"
	"fisherman/infrastructure/log"
	"fisherman/internal"
	"fisherman/internal/handling"
	"fisherman/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(
		map[string]handling.Handler{},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "test"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "'test' is not valid hook name")
}

func TestCommand_Run(t *testing.T) {
	command := handle.NewCommand(
		map[string]handling.Handler{
			"pre-commit": mocks.NewHandlerMock(t).
				IsConfiguredMock.Return(false),
		},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.NoError(t, err)
}

func TestCommand_Run_Hander(t *testing.T) {
	command := handle.NewCommand(
		map[string]handling.Handler{
			"pre-commit": mocks.NewHandlerMock(t).
				IsConfiguredMock.Return(true).
				HandleMock.Return(errors.New("test error")),
		},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "test error")
}

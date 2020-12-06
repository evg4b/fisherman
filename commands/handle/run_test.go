package handle_test

import (
	"errors"
	"fisherman/commands/handle"
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/hookfactory"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(
		hookfactory.HandlerList{},
		&configuration.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "test"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "'test' is not valid hook name")
}

func TestCommand_Run(t *testing.T) {
	command := handle.NewCommand(
		hookfactory.HandlerList{
			"pre-commit": hookfactory.HandlerRegistration{
				Handler:    mocks.NewHandlerMock(t),
				Registered: false,
			},
		},
		&configuration.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.NoError(t, err)
}

func TestCommand_Run_Hander(t *testing.T) {
	command := handle.NewCommand(
		hookfactory.HandlerList{
			"pre-commit": hookfactory.HandlerRegistration{
				Handler: mocks.NewHandlerMock(t).
					HandleMock.Return(errors.New("test error")),
				Registered: true,
			},
		},
		&configuration.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "test error")
}

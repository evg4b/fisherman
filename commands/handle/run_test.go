package handle_test

import (
	"errors"
	"fisherman/commands/handle"
	"fisherman/commands/handle/hooks"
	"fisherman/config"
	"fisherman/internal"
	"fisherman/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(
		hooks.HandlerList{},
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
		hooks.HandlerList{
			"pre-commit": hooks.HandlerRegistration{
				Handler:    mocks.NewHandlerMock(t),
				Registered: false,
			},
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
		hooks.HandlerList{
			"pre-commit": hooks.HandlerRegistration{
				Handler: mocks.NewHandlerMock(t).
					HandleMock.Return(errors.New("test error")),
				Registered: true,
			},
		},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "test error")
}

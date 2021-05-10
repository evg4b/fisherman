package handle_test

import (
	"errors"
	"fisherman/commands/handle"
	"fisherman/infrastructure/log"
	"fisherman/testing/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("test").Return(nil, errors.New("'test' is not valid hook name")),
		mocks.NewCtxFactoryMock(t),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "test"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "'test' is not valid hook name")
}

func TestCommand_Run(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("pre-commit").
			Return(mocks.NewHandlerMock(t).HandleMock.Return(nil), nil),
		mocks.NewCtxFactoryMock(t),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.NoError(t, err)
}

func TestCommand_Run_Hander(t *testing.T) {
	handler := mocks.NewHandlerMock(t).
		HandleMock.Return(errors.New("test error"))
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("pre-commit").Return(handler, nil),
		mocks.NewCtxFactoryMock(t),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run()

	assert.Error(t, err, "test error")
}

package handle_test

import (
	"errors"
	"fisherman/commands/handle"
	"fisherman/infrastructure"
	"fisherman/internal/appcontext"
	"fisherman/pkg/log"
	"fisherman/testing/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var globalVars = map[string]interface{}{
	"BranchName": "/refs/head/develop",
	"Tag":        "1.0.0",
	"UserEmail":  "evg4b@mail.com",
	"UserName":   "evg4b",
}

func getCtx(t *testing.T) *appcontext.ApplicationContext {
	return appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithRepository(mocks.NewRepositoryMock(t).
			GetCurrentBranchMock.Return("/refs/head/develop", nil).
			GetLastTagMock.Return("1.0.0", nil).
			GetUserMock.Return(infrastructure.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
		).
		WithShell(mocks.NewShellMock(t)).
		Build()
}

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("test", globalVars).Return(nil, errors.New("'test' is not valid hook name")),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "test"})
	assert.NoError(t, err)

	err = command.Run(getCtx(t))

	assert.Error(t, err, "'test' is not valid hook name")
}

func TestCommand_Run(t *testing.T) {
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("pre-commit", globalVars).
			Return(mocks.NewHandlerMock(t).HandleMock.Return(nil), nil),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run(getCtx(t))

	assert.NoError(t, err)
}

func TestCommand_Run_Hander(t *testing.T) {
	handler := mocks.NewHandlerMock(t).
		HandleMock.Return(errors.New("test error"))
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("pre-commit", globalVars).Return(handler, nil),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run(getCtx(t))

	assert.Error(t, err, "test error")
}

func TestCommand_Run_GlobalVarsGettingFail(t *testing.T) {
	handler := mocks.NewHandlerMock(t).
		HandleMock.Return(nil)
	command := handle.NewCommand(
		mocks.NewFactoryMock(t).
			GetHookMock.Expect("pre-commit", globalVars).Return(handler, nil),
		&mocks.HooksConfigStub,
		mocks.AppInfoStub,
	)

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	ctx := appcontext.NewContextBuilder().
		WithFileSystem(mocks.NewFileSystemMock(t)).
		WithRepository(mocks.NewRepositoryMock(t).
			GetCurrentBranchMock.Return("/refs/head/develop", nil).
			GetLastTagMock.Return("1.0.0", errors.New("test error")).
			GetUserMock.Return(infrastructure.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
		).
		WithShell(mocks.NewShellMock(t)).
		Build()

	err = command.Run(ctx)

	assert.Error(t, err, "test error")
}

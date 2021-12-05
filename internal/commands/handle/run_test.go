package handle_test

import (
	"errors"
	"fisherman/internal/appcontext"
	. "fisherman/internal/commands/handle"
	"fisherman/internal/constants"
	"fisherman/pkg/vcs"
	"fisherman/testing/mocks"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var globalVars = map[string]interface{}{
	constants.BranchNameVariable:       "/refs/head/develop",
	constants.TagVariable:              "1.0.0",
	constants.UserEmailVariable:        "evg4b@mail.com",
	constants.UserNameVariable:         "evg4b",
	constants.FishermanVersionVariable: constants.Version,
	constants.CwdVariable:              "~/project",
	constants.OsVariable:               runtime.GOOS,
}

func TestCommand_Run(t *testing.T) {
	t.Run("runs correctly", func(t *testing.T) {
		command := NewCommand(
			mocks.NewFactoryMock(t).
				GetHookMock.Expect("pre-commit", globalVars).
				Return(mocks.NewHandlerMock(t).HandleMock.Return(nil), nil),
			&mocks.HooksConfigStub,
			mocks.AppInfoStub,
		)

		err := command.Init([]string{"--hook", "pre-commit"})
		assert.NoError(t, err)

		err = command.Run(appContextStub(t))

		assert.NoError(t, err)
	})

	t.Run("unknown hook", func(t *testing.T) {
		command := NewCommand(
			mocks.NewFactoryMock(t).
				GetHookMock.Expect("test", globalVars).Return(nil, errors.New("'test' is not valid hook name")),
			&mocks.HooksConfigStub,
			mocks.AppInfoStub,
		)

		err := command.Init([]string{"--hook", "test"})
		assert.NoError(t, err)

		err = command.Run(appContextStub(t))

		assert.EqualError(t, err, "'test' is not valid hook name")
	})

	t.Run("call handler and return error", func(t *testing.T) {
		handler := mocks.NewHandlerMock(t).
			HandleMock.Return(errors.New("test error"))

		command := NewCommand(
			mocks.NewFactoryMock(t).
				GetHookMock.Expect("pre-commit", globalVars).Return(handler, nil),
			&mocks.HooksConfigStub,
			mocks.AppInfoStub,
		)

		err := command.Init([]string{"--hook", "pre-commit"})
		assert.NoError(t, err)

		err = command.Run(appContextStub(t))

		assert.EqualError(t, err, "test error")
	})

	t.Run("call handler with global variables", func(t *testing.T) {
		handler := mocks.NewHandlerMock(t).
			HandleMock.Return(nil)

		command := NewCommand(
			mocks.NewFactoryMock(t).
				GetHookMock.Expect("pre-commit", globalVars).Return(handler, nil),
			&mocks.HooksConfigStub,
			mocks.AppInfoStub,
		)

		err := command.Init([]string{"--hook", "pre-commit"})
		assert.NoError(t, err)

		ctx := appcontext.NewContext(
			appcontext.WithFileSystem(mocks.NewFilesystemMock(t)),
			appcontext.WithRepository(mocks.NewRepositoryMock(t).
				GetCurrentBranchMock.Return("/refs/head/develop", nil).
				GetLastTagMock.Return("1.0.0", errors.New("test error")).
				GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
			),
		)

		err = command.Run(ctx)

		assert.EqualError(t, err, "test error")
	})
}

func appContextStub(t *testing.T) *appcontext.ApplicationContext {
	t.Helper()

	return appcontext.NewContext(
		appcontext.WithFileSystem(mocks.NewFilesystemMock(t)),
		appcontext.WithRepository(mocks.NewRepositoryMock(t).
			GetCurrentBranchMock.Return("/refs/head/develop", nil).
			GetLastTagMock.Return("1.0.0", nil).
			GetUserMock.Return(vcs.User{UserName: "evg4b", Email: "evg4b@mail.com"}, nil),
		),
		appcontext.WithCwd("~/project"),
	)
}

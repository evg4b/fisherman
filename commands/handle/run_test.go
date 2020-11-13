package handle_test

import (
	"context"
	"errors"
	"fisherman/clicontext"
	"fisherman/commands/handle"
	"fisherman/config"
	"fisherman/handlers"
	"fisherman/infrastructure/log"
	"fisherman/mocks"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCommand_Run_UnknownHook(t *testing.T) {
	command := handle.NewCommand(flag.ExitOnError, map[string]handlers.Handler{})

	err := command.Init([]string{"--hook", "test"})
	assert.NoError(t, err)

	err = command.Run(ctx(t))

	assert.Error(t, err, "'test' is not valid hook name")
}

func TestCommand_Run(t *testing.T) {
	command := handle.NewCommand(flag.ExitOnError, map[string]handlers.Handler{
		"pre-commit": mocks.NewHandlerMock(t).
			IsConfiguredMock.Return(false),
	})

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run(ctx(t))

	assert.NoError(t, err)
}

func TestCommand_Run_Hander(t *testing.T) {
	c := ctx(t)

	command := handle.NewCommand(flag.ExitOnError, map[string]handlers.Handler{
		"pre-commit": mocks.NewHandlerMock(t).
			IsConfiguredMock.Return(true).
			HandleMock.Return(errors.New("test error")),
	})

	err := command.Init([]string{"--hook", "pre-commit"})
	assert.NoError(t, err)

	err = command.Run(c)

	assert.Error(t, err, "test error")
}

func ctx(t *testing.T) *clicontext.CommandContext {
	return clicontext.NewContext(context.TODO(), clicontext.Args{
		FileSystem: mocks.NewFileSystemMock(t),
		Shell:      mocks.NewShellMock(t),
		Config:     &config.DefaultConfig,
		App:        &clicontext.AppInfo{},
	})
}

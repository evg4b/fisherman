package runner_test

import (
	"context"
	"errors"
	"fisherman/commands"
	"fisherman/config"
	"fisherman/infrastructure"
	"fisherman/mocks"

	"fisherman/runner"
	"io/ioutil"
	"log"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_Run(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		name          string
		args          []string
		commands      []commands.CliCommand
		expectedError error
	}{
		{
			name: "Should run called commnad and return its error",
			args: []string{"init"},
			commands: []commands.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeExpectedCommand(t, "init", errors.New("expected error")),
			},
			expectedError: errors.New("expected error"),
		},
		{
			name: "Should run called commnad and return nil when command executed witout error",
			args: []string{"init"},
			commands: []commands.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeExpectedCommand(t, "init", nil),
			},
			expectedError: nil,
		},
		{
			name: "Should return error when command not found",
			args: []string{"not"},
			commands: []commands.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeCommand(t, "init"),
			},
			expectedError: errors.New("unknown command: not"),
		},
		{
			name:          "Should return error when command not registered",
			args:          []string{"not"},
			commands:      []commands.CliCommand{},
			expectedError: errors.New("unknown command: not"),
		},
		{
			name: "Should not return error when commnad not specified",
			args: []string{},
			commands: []commands.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeCommand(t, "init"),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runnerInstance := runner.NewRunner(context.TODO(), runner.Args{
				Commands: tt.commands,
				Config: &config.FishermanConfig{
					GlobalVariables: make(map[string]interface{}),
				},
				ConfigInfo: &config.LoadInfo{},
				Cwd:        "demo",
				Files:      mocks.NewFileSystemMock(t),
				SystemUser: &user.User{},
				Executable: "bin",
				Repository: mocks.NewRepositoryMock(t).
					GetUserMock.Return(infrastructure.User{}, nil),
			})

			assert.NotPanics(t, func() {
				err := runnerInstance.Run(tt.args)
				assert.Equal(t, tt.expectedError, err)
			})
		})
	}
}

func makeCommand(t *testing.T, name string) *mocks.CliCommandMock {
	return mocks.NewCliCommandMock(t).
		NameMock.Return(name).
		InitMock.Return(nil)
}

func makeExpectedCommand(t *testing.T, name string, err error) *mocks.CliCommandMock {
	return makeCommand(t, name).
		RunMock.Return(err)
}

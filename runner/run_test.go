package runner_test

import (
	"errors"
	"fisherman/commands"
	"fisherman/config"
	"fisherman/infrastructure"
	"fisherman/infrastructure/io"
	commandsmock "fisherman/mocks/commands"
	mocks "fisherman/mocks/infrastructure"
	"fisherman/runner"
	"io/ioutil"
	"log"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunner_Run(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	repo := mocks.Repository{}
	repo.On("GetUser").Return(infrastructure.User{}, nil)

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
				makeCommand("handle"),
				makeCommand("remove"),
				makeExpectedCommand("init", errors.New("expected error")),
			},
			expectedError: errors.New("expected error"),
		},
		{
			name: "Should run called commnad and return nil when command executed witout error",
			args: []string{"init"},
			commands: []commands.CliCommand{
				makeCommand("handle"),
				makeCommand("remove"),
				makeExpectedCommand("init", nil),
			},
			expectedError: nil,
		},
		{
			name: "Should return error when command not found",
			args: []string{"not"},
			commands: []commands.CliCommand{
				makeCommand("handle"),
				makeCommand("remove"),
				makeCommand("init"),
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
				makeCommand("handle"),
				makeCommand("remove"),
				makeCommand("init"),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runnerInstance := runner.NewRunner(runner.NewRunnerArgs{
				CommandList: tt.commands,
				Config: &config.FishermanConfig{
					GlobalVariables: make(map[string]interface{}),
				},
				ConfigInfo: &config.LoadInfo{},
				Cwd:        "demo",
				Files:      &io.LocalFileAccessor{},
				SystemUser: &user.User{},
				Executable: "bin",
				Repository: &repo,
			})

			assert.NotPanics(t, func() {
				err := runnerInstance.Run(tt.args)
				assert.Equal(t, tt.expectedError, err)
			})
		})
	}
}

func makeCommand(name string) *commandsmock.CliCommand {
	command := commandsmock.CliCommand{}
	command.On("Name").Return(name)
	command.On("Init", mock.Anything).Return(nil)

	return &command
}

func makeExpectedCommand(name string, err error) *commandsmock.CliCommand {
	command := makeCommand(name)
	command.On("Run", mock.Anything, mock.Anything).Return(err)

	return command
}

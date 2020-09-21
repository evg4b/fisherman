package runner_test

import (
	"errors"
	"fisherman/commands"
	"fisherman/config"
	"fisherman/infrastructure/io"
	commandsmock "fisherman/mocks/commands"
	loggermock "fisherman/mocks/infrastructure/logger"
	"fisherman/runner"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunner_Run(t *testing.T) {
	logger := loggermock.Logger{}
	logger.On("Debugf", mock.Anything)
	logger.On("Debug", mock.Anything)
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("Write", mock.Anything).Return(1, nil)

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
				makeExpectedCommand("init", errors.New("Expected error")),
			},
			expectedError: errors.New("Expected error"),
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
			expectedError: errors.New("Unknown command: not"),
		},
		{
			name:          "Should return error when command not registered",
			args:          []string{"not"},
			commands:      []commands.CliCommand{},
			expectedError: errors.New("Unknown command: not"),
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
				Config:      &config.DefaultConfig,
				Logger:      &logger,
				ConfigInfo:  &config.ConfigInfo{},
				Cwd:         "demo",
				Files:       &io.LocalFileAccessor{},
				SystemUser:  &user.User{},
				Executable:  "bin",
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
	return &command
}

func makeExpectedCommand(name string, err error) *commandsmock.CliCommand {
	command := makeCommand(name)
	command.On("Name").Return(name)
	command.On("Run", mock.Anything, mock.Anything).Return(err)
	return command
}

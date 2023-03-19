package app_test

import (
	"context"
	"fisherman/internal"
	. "fisherman/internal/app"
	"fisherman/internal/constants"
	"fisherman/pkg/log"
	"fisherman/testing/mocks"
	"fisherman/testing/testutils"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/go-errors/errors"

	"github.com/stretchr/testify/assert"
)

func TestRunner_Run(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		name        string
		args        []string
		commands    []internal.CliCommand
		expectedErr string
	}{
		{
			name: "should run called commnad and return its error",
			args: []string{"init"},
			commands: []internal.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeExpectedCommand(t, "init", errors.New("expected error")),
			},
			expectedErr: "expected error",
		},
		{
			name: "should run called commnad and return nil when command executed witout error",
			args: []string{"init"},
			commands: []internal.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeExpectedCommand(t, "init", nil),
			},
		},
		{
			name: "should return error when command not found",
			args: []string{"not"},
			commands: []internal.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeCommand(t, "init"),
			},
			expectedErr: "unknown command: not",
		},
		{
			name:        "should return error when command not registered",
			args:        []string{"not"},
			commands:    []internal.CliCommand{},
			expectedErr: "unknown command: not",
		},
		{
			name: "should not return error when commnad not specified",
			args: []string{},
			commands: []internal.CliCommand{
				makeCommand(t, "handle"),
				makeCommand(t, "remove"),
				makeCommand(t, "init"),
			},
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			app := NewFishermanApp(
				WithCommands(testCase.commands),
				WithCwd("/"),
			)

			assert.NotPanics(t, func() {
				err := app.Run(context.TODO(), testCase.args)
				testutils.AssertError(t, testCase.expectedErr, err)
			})
		})
	}

	t.Run("default print", func(t *testing.T) {
		buffer := strings.Builder{}
		log.SetOutput(&buffer)
		defer log.SetOutput(io.Discard)

		expectedCommand1 := "init"
		expectedCommand2 := "version"
		expectedCommand3 := "handle"
		expectedCommand4 := "remove"

		app := NewFishermanApp(
			WithCommands([]internal.CliCommand{
				makeCommand(t, expectedCommand1),
				makeCommand(t, expectedCommand2),
				makeCommand(t, expectedCommand3),
				makeCommand(t, expectedCommand4),
			}),
			WithCwd("/"),
		)

		err := app.Run(context.TODO(), []string{})

		output := buffer.String()

		assert.NoError(t, err)
		assert.Contains(t, output, constants.AppName)
		assert.Contains(t, output, constants.Version)
		assert.Contains(t, output, expectedCommand1)
		assert.Contains(t, output, expectedCommand2)
		assert.Contains(t, output, expectedCommand3)
		assert.Contains(t, output, expectedCommand4)
		assert.Contains(t, output, "Small git hook management tool for developer.")
	})
}

func TestRunner_Interrupt(t *testing.T) {
	chanel := make(chan os.Signal, 1)
	chanel <- os.Interrupt

	commandMock := mocks.NewCliCommandMock(t).
		NameMock.Return("test-command").
		RunMock.Set(func(ctx context.Context, _ []string) error {
		<-ctx.Done()

		return ctx.Err()
	})

	app := NewFishermanApp(
		WithCommands([]internal.CliCommand{commandMock}),
		WithCwd("/"),
		WithInterruptChanel(chanel),
	)

	err := app.Run(context.Background(), []string{"test-command"})

	assert.EqualError(t, err, context.Canceled.Error())
}

func makeCommand(t *testing.T, name string) *mocks.CliCommandMock {
	t.Helper()

	return mocks.NewCliCommandMock(t).
		NameMock.Return(name).
		DescriptionMock.Return(fmt.Sprintf("This is %s command", name))
}

func makeExpectedCommand(t *testing.T, name string, err error) *mocks.CliCommandMock {
	t.Helper()

	return makeCommand(t, name).
		RunMock.Return(err)
}

package version_test

import (
	"bytes"
	"fisherman/commands/version"
	"fisherman/constants"
	"fisherman/infrastructure/log"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	output := bytes.NewBufferString("")
	log.SetOutput(output)

	command := version.NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)

	err = command.Run()

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintln(constants.Version), output.String())
}

func TestCommand_Description(t *testing.T) {
	command := version.NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)
	assert.NotEmpty(t, command.Description())
}

func TestCommand_Name(t *testing.T) {
	command := version.NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)
	assert.Equal(t, "version", command.Name())
}
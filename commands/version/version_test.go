package version_test

import (
	"bytes"
	"fisherman/commands/version"
	"fisherman/infrastructure/log"
	"fisherman/testing/mocks"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCommand_Run(t *testing.T) {
	output := bytes.NewBufferString("")
	log.SetOutput(output)

	command := version.NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)

	err = command.Run(mocks.NewExecutionContextMock(t))

	assert.NoError(t, err)
	assert.Equal(t, "fisherman@x.x.x", output.String())
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

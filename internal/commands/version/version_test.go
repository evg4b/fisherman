package version_test

import (
	"bytes"
	"context"
	. "fisherman/internal/commands/version"
	"fisherman/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Run(t *testing.T) {
	output := bytes.NewBufferString("")
	log.SetOutput(output)

	command := NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)

	err = command.Run(context.TODO())

	assert.NoError(t, err)
	assert.Equal(t, "fisherman@x.x.x", output.String())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)
	assert.NotEmpty(t, command.Description())
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand()
	err := command.Init([]string{})

	assert.NoError(t, err)
	assert.Equal(t, "version", command.Name())
}

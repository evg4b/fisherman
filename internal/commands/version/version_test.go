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

	err := command.Run(context.TODO(), []string{})

	assert.NoError(t, err)
	assert.Equal(t, "fisherman@x.x.x\n", output.String())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand()

	assert.NotEmpty(t, command.Description())
}

func TestCommand_Name(t *testing.T) {
	command := NewCommand()

	assert.Equal(t, "version", command.Name())
}

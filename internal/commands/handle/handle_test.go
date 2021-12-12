package handle_test

import (
	. "fisherman/internal/commands/handle"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := NewCommand()

	assert.Equal(t, "handle", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := NewCommand()

	assert.NotEmpty(t, command.Description())
}

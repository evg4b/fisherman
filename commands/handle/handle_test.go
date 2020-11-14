package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/internal/handling"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(flag.ExitOnError, make(map[string]handling.Handler))

	assert.Equal(t, "handle", command.Name())
}

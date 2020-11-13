package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/handlers"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(flag.ExitOnError, make(map[string]handlers.Handler))

	assert.Equal(t, "handle", command.Name())
}

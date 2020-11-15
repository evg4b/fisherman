package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/config"
	"fisherman/internal"
	"fisherman/internal/handling"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(
		make(map[string]handling.Handler),
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	assert.Equal(t, "handle", command.Name())
}

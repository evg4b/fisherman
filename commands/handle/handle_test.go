package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/config"
	"fisherman/internal"
	"fisherman/internal/hookfactory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(
		hookfactory.HandlerList{},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	assert.Equal(t, "handle", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := handle.NewCommand(
		hookfactory.HandlerList{},
		&config.HooksConfig{},
		&internal.AppInfo{},
	)

	assert.NotEmpty(t, command.Description())
}

package handle_test

import (
	"fisherman/commands/handle"
	"fisherman/configuration"
	"fisherman/internal"
	"fisherman/internal/hookfactory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Name(t *testing.T) {
	command := handle.NewCommand(
		&hookfactory.TFactory{},
		&configuration.HooksConfig{},
		&internal.AppInfo{},
	)

	assert.Equal(t, "handle", command.Name())
}

func TestCommand_Description(t *testing.T) {
	command := handle.NewCommand(
		&hookfactory.TFactory{},
		&configuration.HooksConfig{},
		&internal.AppInfo{},
	)

	assert.NotEmpty(t, command.Description())
}

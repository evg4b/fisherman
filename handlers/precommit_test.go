package handlers_test

import (
	"fisherman/commands"
	"fisherman/config"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCommitHandler(t *testing.T) {
	assert.NotPanics(t, func() {
		err := handlers.PreCommitHandler(&commands.CommandContext{
			Config: &config.HooksConfig{},
		}, []string{})
		assert.NoError(t, err)
	})
}

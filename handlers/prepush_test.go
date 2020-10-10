package handlers_test

import (
	"fisherman/commands"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrePushHandler(t *testing.T) {
	assert.NotPanics(t, func() {
		err := handlers.PrePushHandler(&commands.CommandContext{}, []string{})
		assert.NoError(t, err)
	})
}

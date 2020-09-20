package handlers_test

import (
	"fisherman/commands/context"
	"fisherman/constants"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreApplyPatchHandler(t *testing.T) {
	constants.Version = "1.0.1"
	assert.NotPanics(t, func() {
		err := handlers.PreApplyPatchHandler(&context.CommandContext{}, []string{})
		assert.Error(t, err, "This hook is not supported in version 1.0.1.")
	})
}

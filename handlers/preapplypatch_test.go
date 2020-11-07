package handlers_test

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreApplyPatchHandler(t *testing.T) {
	constants.Version = testVertion
	assert.NotPanics(t, func() {
		err := handlers.PreApplyPatchHandler(&clicontext.CommandContext{}, []string{})
		assert.Error(t, err, "This hook is not supported in version 1.0.1.")
	})
}

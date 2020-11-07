package handlers_test

import (
	"fisherman/clicontext"
	"fisherman/constants"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVertion = "1.0.1"

func TestApplyPatchMsgHandler(t *testing.T) {
	constants.Version = testVertion
	assert.NotPanics(t, func() {
		err := handlers.ApplyPatchMsgHandler(&clicontext.CommandContext{}, []string{})
		assert.Error(t, err, "This hook is not supported in version 1.0.1.")
	})
}

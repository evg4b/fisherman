package handlers_test

import (
	"fisherman/clicontext"
	"fisherman/config"
	"fisherman/constants"
	"fisherman/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVertion = "1.0.1"

func TestNotSupportedHandler(t *testing.T) {
	constants.Version = testVertion
	assert.NotPanics(t, func() {
		err := new(handlers.NotSupportedHandler).Handle(&clicontext.CommandContext{}, []string{})
		assert.Error(t, err, "This hook is not supported in version 1.0.1.")
	})
}

func TestNotSupportedHandler_IsConfigured(t *testing.T) {
	var handler handlers.NotSupportedHandler
	tests := []struct {
		name   string
		config *config.HooksConfig
	}{
		{name: "empty config", config: &config.HooksConfig{}},
		{name: "not empty config", config: &config.DefaultConfig.Hooks},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, handler.IsConfigured(tt.config))
		})
	}
}

package handling_test

import (
	"fisherman/config"
	"fisherman/constants"
	"fisherman/internal/handling"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVertion = "1.0.1"

func TestNotSupportedHandler(t *testing.T) {
	constants.Version = testVertion
	assert.NotPanics(t, func() {
		err := new(handling.NotSupportedHandler).Handle([]string{})
		assert.Error(t, err, "This hook is not supported in version 1.0.1.")
	})
}

func TestNotSupportedHandler_IsConfigured(t *testing.T) {
	var handler handling.NotSupportedHandler
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

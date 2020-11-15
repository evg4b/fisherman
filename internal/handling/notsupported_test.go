package handling_test

import (
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

package guards_test

import (
	"errors"
	"fisherman/pkg/guards"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoErrorPanicForError(t *testing.T) {
	err := errors.New("Test err")
	assert.PanicsWithError(t, err.Error(), func() {
		guards.NoError(err)
	})
}

func TestNoErrorNotPanicForNil(t *testing.T) {
	assert.NotPanics(t, func() {
		guards.NoError(nil)
	})
}

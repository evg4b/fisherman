package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCriticalErrorPanicForError(t *testing.T) {
	err := errors.New("Test err")
	assert.PanicsWithError(t, err.Error(), func() {
		HandleCriticalError(err)
	})
}

func TestHandleCriticalErrorNotPanicForNil(t *testing.T) {
	assert.NotPanics(t, func() {
		HandleCriticalError(nil)
	})
}

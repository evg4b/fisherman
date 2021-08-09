package utils_test

import (
	"errors"
	"fisherman/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCriticalErrorPanicForError(t *testing.T) {
	err := errors.New("Test err")
	assert.PanicsWithError(t, err.Error(), func() {
		utils.HandleCriticalError(err)
	})
}

func TestHandleCriticalErrorNotPanicForNil(t *testing.T) {
	assert.NotPanics(t, func() {
		utils.HandleCriticalError(nil)
	})
}
